from ctypes import *
import time
import os

word2vec_dll = cdll.LoadLibrary("word2vec.dll")


def clock(func):
    def handle(*args):
        time_start = time.time()
        r = func(*args)
        time_end = time.time()
        print(f"time use: {time_end - time_start}")
        return r
    return handle

@clock
def word2vec_go(train_file, vector_file, window: int, dimention: int, negative: int, lr, show: int):
    """

    :param train_file:
    :param vector_file:
    :param window:
    :param dimention:
    :param negative:
    :param lr:
    :param show: 1 直接调用向量文件、0需要重新训练
    :return:
    """
    if not os.path.exists(vector_file):
        with open(vector_file, "w", encoding="utf-8") as f:
            pass

    word2vec_dll.Run(bytes(train_file, encoding="utf-8"), bytes(vector_file, encoding="utf-8"),
                            window, dimention, negative, c_float(lr), show)


if __name__ == "__main__":
    word2vec_go("data/text8split", "data/vectors", 10, 100, 6, 0.05, 0)