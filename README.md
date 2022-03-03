Bloom filters
-------------

裁减了原仓库的核心代码，增加了 redis bitset 支持，删除了原代码的某些方法，引入 mini redis 保留了原有的一些单元测试代码，移除或注释了目前不支持的测试代码。
目前基础的 Test 和 Add 都已支持，也尽量采用管道批量操作，目前还缺失一个 Count 方法统计总共插入的数量。根据布隆过滤器的原理，如果插入数量超过预设，准确率会
下降，所以要避免这种情况必须知道插入的元素数量是否已经超过了，不过因为目前我自己没有这个需求，同时也是破坏性的修改源码，没打算提 pr，暂时就这样了。

## Design

A Bloom filter has two parameters: _m_, the number of bits used in storage, and _k_, the number of hashing functions on elements of the set. (The actual hashing functions are important, too, but this is not a parameter for this implementation). A Bloom filter is backed by a [BitSet](https://github.com/bits-and-blooms/bitset); a key is represented in the filter by setting the bits at each value of the  hashing functions (modulo _m_). Set membership is done by _testing_ whether the bits at each value of the hashing functions (again, modulo _m_) are set. If so, the item is in the set. If the item is actually in the set, a Bloom filter will never fail (the true positive rate is 1.0); but it is susceptible to false positives. The art is to choose _k_ and _m_ correctly.

In this implementation, the hashing functions used is [murmurhash](https://github.com/spaolacci/murmur3), a non-cryptographic hashing function.


Given the particular hashing scheme, it's best to be empirical about this. Note
that estimating the FP rate will clear the Bloom filter.
