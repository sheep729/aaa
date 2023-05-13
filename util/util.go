package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASdfGHJKLZXCVBNMQWERTYUIOP") //这段代码定义了一个名为 letters 的变量，其类型为 []byte，也就是一个字节切片。
	result := make([]byte, n)                                                    //使用 Go 语言中的 make 函数创建了一个名为 result 的切片，类型为 []byte，n为元素个数，指定切片的长度

	rand.Seed(time.Now().Unix()) //使用 Go 标准库中的 time.Now().Unix() 函数获取当前时间戳，并将其作为随机数生成器的种子。然后，rand.Seed() 函数使用该种子来初始化默认的全局随机数生成器。
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	} //使用 for 循环遍历 result 切片中的每个元素，并为其随机生成一个字符
	//使用 rand.Intn(n) 函数从 letters 切片中随机选择一个字符，并将其赋值给 result 切片中对应的元素。
	return string(result)
}
