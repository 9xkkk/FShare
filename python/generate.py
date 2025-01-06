# -*- coding: utf-8 -*-
import hashlib
import sys


def generate(sp_id):
    # 参数已写定
    secretKey = "nurseryDatabase"

    # 生成指纹
    sp_id = str(sp_id)  # 转换为字符串
    rand_str = secretKey + sp_id
    fp = hashlib.md5(rand_str.encode('gbk', errors='ignore')).hexdigest()  # 十六进制
    fp = bin(int(fp, 16))[2:].zfill(128)  # 转二进制，舍弃0b前缀，补充前导0至128位

    return fp


if __name__ == '__main__':
    fp = generate(sys.argv[1])
    if fp == "":
        print("false")
    else:
        print(fp)
