百度云加密文件密码爆破工具
==========================

下载地址
--------

百度盘: http://pan.baidu.com/s/1c0FeVgk

备用: http://pan.baidu.com/s/1nty78sL 密码: vvsb

使用方法
--------

下载软件输入[命令行](#windows下如何快速输入命令行)

`bdptester -u "http://pan.baidu.com/share/init?shareid=3521213983&uk=3793282542"`

接下来就是耐心等待

注:

+ 根据__网速__的不同 一般__100M__的网每秒可以测试__3000__条左右
+ 而所有密码的可能性有__1679616__种
+ 如果测试到一半就对了 那么就是要花__1679616/2/3000=279.936秒__的时间

还是挺快的不是吗! enjoy it~

断点续试/测试可用性的方法
--------

如果知道密码可以用`-at`参数来试试速度

比如密码是`vvsb` 就输入命令行

`bdptester -u "http://pan.baidu.com/share/init?shareid=3521213983&uk=3793282542" -at v000`

就是从`v000`开始测试 基本上几秒钟就出结果了

+ 这个参数也适用于上次测了一半中途停了, 因为程序会输出现在测到哪里 所以下一次再测试就可以直接填at参数而不用从头开始了

手动设定线程数
--------------

默认的线程数是500

你也可以使用-j参数来改变他

`bdptester -u "something" -j 1000`

即为使用1000个线程并行

注: 

+ 这里的线程并不是指操作系统线程 资源占用并不是很多 当然跑快了CPU负荷还是挺高的233
+ 所以稍微多点其实也没啥关系
+ 不过一般500就足够跑满网速没必要再多了
+ 当然如果你是1000M的土豪可以试试调高...

一般网速和线程数的对应
---------------------

| 网速       | 建议线程数   |
| ---------- | ------------ |
| 4M         | 20           |
| 10M        | 50           |
| 20M        | 100          |
| 50M        | 200          |
| 100M       | 400          |


Windows下如何快速输入命令行?
----------------------------

其实不用win+r然后cmd也可以很方便的输入命令

只要在资源管理器中打开程序所在目录 在空白处按shift键+鼠标右键 就可以出现含有"在此处打开命令窗口的界面"

然后你可以在里面打命令 也可以在外面写好 然后用右键粘贴进去就可以运行了