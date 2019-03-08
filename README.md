# word2vec_python_go
go实现的word2vec,使用负采样算法。打包成word2vec.dll后供python调用<br>
go代码编译命令:go build -buildmode=c-shared -o word2vec.dll  src/main/word2vec.go<br>

python测试代码如下<br>

from ctypes import *<br>
import time<br>
import os<br>

word2vec_dll = cdll.LoadLibrary("word2vec.dll")<br>


def clock(func):<br>
    def handle(*args):<br>
        time_start = time.time()<br>
        r = func(*args)<br>
        time_end = time.time()<br>
        print(f"time use: {time_end - time_start}")<br>
        return r<br>
    return handle<br>

@clock<br>
def word2vec_go(train_file, vector_file, window: int, dimention: int, negative: int, lr, show: int):<br>
    """<br>

    :param train_file:<br>
    :param vector_file:<br>
    :param window:<br>
    :param dimention:<br>
    :param negative:<br>
    :param lr:<br>
    :param show: 1 直接调用向量文件、0需要重新训练<br>
    :return:<br>
    """<br>
    if not os.path.exists(vector_file):<br>
        with open(vector_file, "w", encoding="utf-8") as f:<br>
            pass<br>

    word2vec_dll.Run(bytes(train_file, encoding="utf-8"), bytes(vector_file, encoding="utf-8"),<br>
                            window, dimention, negative, c_float(lr), show)<br>


if __name__ == "__main__":<br>
    word2vec_go("data/text8split", "data/vectors", 10, 100, 6, 0.05, 0)<br>