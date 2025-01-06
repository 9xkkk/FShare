# -*- coding: utf-8 -*-
import sys
import numpy as np
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
from sklearn.ensemble import RandomForestClassifier
from sklearn.linear_model import LogisticRegression
from sklearn.svm import SVC
from sklearn.neighbors import KNeighborsClassifier
from sklearn.tree import DecisionTreeClassifier
import joblib

def train_analysis(Original_path):
    # 定义分类器列表（去掉神经网络）
    names = [
        "Random Forest",
        "Logistic Regression",
        "Linear SVM",
        "RBF SVM",
        "K-Nearest Neighbors",
        "Decision Tree",
    ]

    classifiers = [
        RandomForestClassifier(random_state=42),
        LogisticRegression(max_iter=1000, random_state=42),
        SVC(kernel="linear", C=0.025, random_state=42),
        SVC(gamma=2, C=1, random_state=42),
        KNeighborsClassifier(3),
        DecisionTreeClassifier(max_depth=5, random_state=42),
    ]

    # 保存精确度的数组
    Original_accuracy = []
    W_Precision = []
    W_Recall = []
    W_F1 = []

    # -----------------train original model ------------------- #
    # 读入原始数据集
    Original_Class_df = pd.read_csv(Original_path)

    # 划分数据集
    Original_Class_label_df = Original_Class_df['Class']  # 以Class列作为预测标签
    Original_Class_feature_df = Original_Class_df.drop(['id', 'Class'], axis=1)  # 特征列
    Original_X_train, X_test, Original_Y_train, Y_test = train_test_split(
        Original_Class_feature_df, Original_Class_label_df, test_size=0.2, random_state=42
    )

    # 训练和评估每个分类器
    for name, clf in zip(names, classifiers):
        clf.fit(Original_X_train, Original_Y_train)

        # 评估模型在测试集上的表现
        accuracy = clf.score(X_test, Y_test)
        Original_accuracy.append(accuracy)  # 保存到 Original_accuracy 数组
        Original_Y_pred = clf.predict(X_test)

        report_dict = classification_report(Y_test, Original_Y_pred, output_dict=True)
        precision = report_dict['weighted avg']['precision']
        recall = report_dict['weighted avg']['recall']
        f1 = report_dict['weighted avg']['f1-score']

        W_Precision.append(precision)
        W_Recall.append(recall)
        W_F1.append(f1)


        # 保存每个模型
        model_filename = f'Train_Models/original_models/{name.replace(" ", "_")}_model.pkl'
        joblib.dump(clf, model_filename)


    return [[Original_accuracy, list(np.zeros(len(Original_accuracy)).tolist())],
            [W_Precision, list(np.zeros(len(W_Precision)).tolist())],
            [W_Recall, list(np.zeros(len(W_Recall)).tolist())],
            [W_F1, list(np.zeros(len(W_F1)).tolist())]]



if __name__ == '__main__':
    result = train_analysis(sys.argv[1])
    print(result)





