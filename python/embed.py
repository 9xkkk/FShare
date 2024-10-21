# -*- coding: utf-8 -*-
import hashlib
import math
import numpy
import pandas
import csv
import random
import sympy
import os
import sys
from scipy.special import roots_legendre


def embed(filename, sp_id):
    # 参数已写定
    sensitivity = 1
    secretKey = "nurseryDatabase"
    epsilon = 3
    omega = 2  # omega取2有错

    # 生成指纹
    sp_id = str(sp_id)  # 转换为字符串
    rand_str = secretKey + sp_id
    fp = hashlib.md5(rand_str.encode('gbk', errors='ignore')).hexdigest()  # 十六进制
    fp = bin(int(fp, 16))[2:].zfill(128)  # 转二进制，舍弃0b前缀，补充前导0至128位

    # 读取文件
    R_DF = pandas.read_csv("./csvfile/"+filename, encoding='utf-8-sig')
    R = numpy.array(R_DF)
    # R = R[0:, 0:9]  # 删除了最后一列class属性

    # 嵌入过程
    rows = numpy.size(R, 0)
    columns = numpy.size(R, 1)
    K = math.floor(math.log(sensitivity, 2)) + 1
    p = 1 / (math.exp(epsilon / K) + 1)
    bits_required = numpy.floor(numpy.log2(numpy.max(R, axis=0))) + 1  # 数据的二进制位数,第0列为主键位数
    L = len(fp)
    T = columns - 1  # 属性数
    for i in range(rows):
        primary_key_att_value = R[i, 0]  # 主键
        for t in range(1, columns):
            bit_length = int(bits_required[t])
            K_min = min(bit_length, K)
            r_it_binary = bin(R[i, t])[2:].zfill(bit_length)  # 数据转为二进制，去0b前缀，补充前导零至位数与该属性的最大位数一致
            ascii_values = [ord(char) for char in secretKey]  # 字符串转ascii
            for k in range(1, K_min + 1):
                seed = sum(ascii_values) + primary_key_att_value + t + k  # numpy.random.seed不接受字符串种子，对种子的每个数求和
                numpy.random.seed(seed)
                rnd_seq = numpy.random.randint(1, numpy.iinfo(numpy.int32).max, size=2)
                if rnd_seq[0] % math.floor(1 / (omega * p)) == 0:
                    l = rnd_seq[1] % L
                    f = fp[l]
                    r_it_binary_list = list(r_it_binary)
                    r_it_binary_list[bit_length - k] = f
                    r_it_binary = ''.join(r_it_binary_list)
            R[i, t] = int(r_it_binary, 2)

    # 建立sp_id与fp的对应关系,追加写入csv文件
    # 该csv第一行无属性名，直接是一条记录，如果01字符串显示不全请调整excel的显示，后期换成数据库应该可以避免该问题
    # 写入方式是追加写入，如果测试时出现错误请清空该表再测试，避免出现意外的bug
    with open("fake_database.csv", "a", encoding="utf-8-sig", newline="") as f:
        csv_writer = csv.writer(f)
        record = [sp_id, fp]  # 第一列是sp_id,第二列是f,一条记录是一个对应关系
        csv_writer.writerow(record)
        f.close()

    # 保存嵌入指纹的数据库
    # 该算法只能处理整数，fmt='%d'代表将所有数据强转为整数形式输出到新表，如果去掉用户则应手动将excel里的数据转换为整数再提取指纹，否则报错
    output_filename = "./csvfile/"+os.path.splitext(filename)[0] + "_FP.csv"
    numpy.savetxt(output_filename, R, delimiter=',', header=','.join(R_DF.columns), fmt='%d', comments='')

    return fp


if __name__ == '__main__':
    fp = embed(sys.argv[1], sys.argv[2])
    if fp == "":
        print("false")
    else:
        print(fp)
