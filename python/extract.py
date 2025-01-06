# -*- coding: utf-8 -*-
"""
readme
请看embed.py
"""
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


def extract(filename, epsilon_avg):
    # 参数已写定,与嵌入需保持一致
    sensitivity = 1
    secretKey = "nurseryDatabase"
    # epsilon = float(epsilon_avg)
    epsilon = 3
    omega = 2  # omega取2有错

    # 读取文件
    MR_DF = pandas.read_csv(filename)
    MR = numpy.array(MR_DF)

    # 提取过程
    rows = numpy.size(MR, 0)
    columns = numpy.size(MR, 1)
    K = math.floor(math.log(sensitivity, 2)) + 1
    p = 1 / (math.exp(epsilon / K) + 1)
    bits_required = numpy.floor(numpy.log2(numpy.max(MR, axis=0))) + 1
    L = 128  # 指纹长度
    T = columns - 1
    fp_count0 = numpy.zeros(128, dtype=int)  # 投票计数
    fp_count1 = numpy.zeros(128, dtype=int)
    for i in range(rows):
        primary_key_att_value = MR[i, 0]
        for t in range(1, columns):
            bit_length = int(bits_required[t])
            K_min = min(bit_length, K)
            r_it_binary = bin(MR[i, t])[2:].zfill(bit_length)
            ascii_values = [ord(char) for char in secretKey]
            for k in range(1, K_min + 1):
                seed = sum(ascii_values) + primary_key_att_value + t + k
                numpy.random.seed(seed)
                rnd_seq = numpy.random.randint(1, numpy.iinfo(numpy.int32).max, size=2)
                if rnd_seq[0] % math.floor(1 / (omega * p)) == 0:
                    l = rnd_seq[1] % L
                    if r_it_binary[bit_length - k] == '0':
                        fp_count0[l] += 1
                    else:
                        fp_count1[l] += 1
    extracted_fp = numpy.where(fp_count1 >= fp_count0, 1, 0)  # 统计指纹的每一位
    extracted_fp = str(extracted_fp)  # 格式转换，下同
    extracted_fp = "".join(extracted_fp.split())  # 去空格
    extracted_fp = extracted_fp.replace("[", "").replace("]", "")  # 去括号

    # 读取数据库找出fp对应的sp_id
    # database_DF = pandas.read_csv("fake_database.csv", header=None)
    # for index, row in database_DF.iterrows():  # 遍历第二列，查找匹配的指纹，返回对应的sp_id
    #     if row.iloc[1] == extracted_fp:
    #         return row.iloc[0], row.iloc[1]
    # return '', extracted_fp  # 如果没有找到匹配的，返回空串
    # database_DF = pandas.read_csv("fake_database.csv", header=None)
    # for index, row in database_DF.iterrows():  # 遍历第二列，查找匹配的指纹，返回对应的sp_id
    #     if row.iloc[1] == extracted_fp:
    #         return row.iloc[0], row.iloc[1]
    return extracted_fp  # 如果没有找到匹配的，返回空串


if __name__ == '__main__':
    FP = extract(sys.argv[1], sys.argv[2])
    if FP == '':
        print("false")
    else:
        print(FP)
