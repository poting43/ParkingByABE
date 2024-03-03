#!/bin/bash

#sed -i -e 's/\r$//' run_multiple_times.sh

g++ time_tbcpabe.cpp bn_pair.cpp zzn12a.cpp ecn2.cpp zzn4.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp miracl.a

for i in {1..10}
do
    ./a.out -> output$i.txt
    echo "Run $i completed."
done

echo "All runs completed."

python average.py
