# 用 miracl 函式庫進行密碼學時間測試

## 1. 環境

linux(建議)

## 2. 檔案說明

1. miracl_library.zip : 為 miracl 函式存放位置
2. time_sign.cpp : 計算簽名與驗證時間
3. time_tbcpabe.cpp : 計算論文中 abe 的正確性與時間

## 3. 使用流程

1. 解壓縮 miracl_library.zip 為 miracl_library 資料夾
2. 將 time_sign.cpp 和 time_tbcpabe.cpp 移至該 miracl_library 資料夾中
3. 用以下指令編譯 .cpp 檔案成 a.out，並且執行
   
``` CMD
g++ time_sign.cpp bn_pair.cpp zzn12a.cpp ecn2.cpp zzn4.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp mrshs256.c miracl.a
./a.out
```

or

``` CMD
g++ time_tbcpabe.cpp bn_pair.cpp zzn12a.cpp ecn2.cpp zzn4.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp mrshs256.c miracl.a
./a.out
```

## 4. 備註

各種實驗中的變數皆可利用調整程式中參數來得到相對應的結果

## 5. 參考資料

[https://blog.csdn.net/qq_44925830/article/details/123657274)
