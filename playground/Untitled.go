package main

import (
	// "encoding/binary"
	// "strconv"
	"fmt"
)

func main() {
	// buf := new(bytes.Buffer)
	data:="1001s"
	var this []byte

	for _,v:= range(data){
		var test byte = v
		this=append(this,test)
	}
	fmt.Println(this)
	// n,_:=buf.WriteString(data)
	// test,err := strconv.ParseInt(data,2,32)
	// binary_val:=strconv.FormatInt(test,2)
	// test2:=[]byte(binary_val)
	// for k,_ := range test2{
	// 	fmt.Println(binary_val[k])
	// 	if(binary_val[k]==48){
	// 		test2[k]=1
	// 	}else{
	// 		test2[k]=0
	// 	}
		

	// }


	// err2 := binary.Write(buf, binary.LittleEndian, test)
	// if err2 != nil {
	// 	fmt.Println("binary.Write failed:", err)
	// }

}
// test,err := strconv.ParseInt(data,2,32)
// 	if err != nil{
// 		fmt.Println("error")
// 	}
// 	s:=fmt.Sprintf("%b",test)
// 	fmt.Println(s)
// 	var binary_val string = strconv.FormatInt(test,2)



// 	var len int =int(math.Ceil(float64(len(binary_val))/2.0))
// 	fmt.Println(len)

// 	var bin_rep [10]byte
// 	var j int = 0
// 	for i:= 0; i<len; i++{
// 		if (len%2 != 0 && i==0){
// 			if binary_val[j]==48{
// 				bin_rep[i]=00
// 			}else{
// 				bin_rep[i]=01
// 			}
// 			j=j+1
// 			continue
// 		}
// 		if (binary_val[j]==48 && binary_val[j+1] ==49){
// 			bin_rep[i]=01
// 		}else if (binary_val[j]==48 && binary_val[j+1]==48){
// 			bin_rep[i]=00
// 		}else if (binary_val[j]==49 && binary_val[j+1]==48){
// 			bin_rep[i]=10
// 		}else{
// 			bin_rep[i]=11
// 		}
// 		j=j+2
		
// 	}
// 	fmt.Println(bin_rep)
// 	test2:=[]byte(binary_val)
// 	for k,_ := range test2{
// 		fmt.Println(binary_val[k])
// 		if(binary_val[k]==48){
// 			test2[k]=00
// 		}else{
// 			test2[k]=01
// 		}
// 	}
// 	return test2