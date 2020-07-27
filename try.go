package main

import(
	"os"
	"fmt"
	"flag"
	"log"
	"strings"
	"strconv"
	"encoding/binary"
	"encoding/hex"
    "golang.org/x/sys/windows/registry"
)

//Function To display the values of the registry Key if Present

func display_values(k registry.Key, value string, full_path string)(custom_error int){

	type_dict := map[uint32]string {0:"NONE",
    1:"SZ",
    2:"EXPAND_SZ", 
    3:"BINARY",
    4:"DWORD",
    5:"DWORD_BIG_ENDIAN",
    6:"LINK",
    7:"MULTI_SZ",
    8:"RESOURCE_LIST",
    9:"FULL_RESOURCE_DESCRIPTOR",
    10:"RESOURCE_REQUIREMENTS_LIST",
	11:"QWORD"}
	
	buf:=make([]byte,1024)
	n,t,err := k.GetValue(value, buf)
	if err == registry.ErrShortBuffer {
		buf=make([]byte,n)
		_,_,err := k.GetValue(value,buf)
		if err != nil{
			log.Println(err)
		}
	}
	if err == registry.ErrNotExist {
		return 1
	}
	
	fmt.Printf("\n\n")
	fmt.Println(full_path+"\n")
	fmt.Printf(value+"\t %v\n",type_dict[t])

	if ( type_dict[t] == "DWORD"){
		test:=binary.LittleEndian.Uint32(buf)
		fmt.Printf("0X%x",test)
	}else if ( type_dict[t] == "DWORD_BIG_ENDIAN"){
		test:=binary.BigEndian.Uint32(buf)
		fmt.Printf("0X%x",test)
	}else if ( type_dict[t] == "BINARY"){
		s:= (buf[:n])
		fmt.Printf("%b",s)
	}else if ( type_dict[t] == "QWORD"){
		s:= binary.LittleEndian.Uint64(buf)
		fmt.Printf("0X%x",s)
	}else{
		s:= string(buf[:n])
		fmt.Printf("%s",string(s))
	}
	return 0
}


//Function to Query the Registry for keys, subkeys and values if given recursively


func query_reg_value(key registry.Key, full_path string, value string, r bool) (custom_error int) {

	
	first := strings.Split(full_path,"\\")
	var relative_path string=(full_path[len(first[0])+1:])

	//opening a key to query the values
	k, err := registry.OpenKey(key, relative_path, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return 1

	}
	defer k.Close()	
	
	

	//If the query is to be recursive
	if r{
		if value==""{
			temp_value,err1 := k.ReadValueNames(-1)
			if err1 == nil{
				for _,val := range temp_value{
					this:=display_values(k,val,full_path)
					fmt.Println(this)

				}
			}
			subkey,err2 := k.ReadSubKeyNames(-1)
			if err2 == nil{
				for _,val := range subkey{
					query_reg_value(key,full_path+"\\"+val,"",r)
					}
			}
			
		}else{
			display_values(k,value,full_path)
			subkey,err2 := k.ReadSubKeyNames(-1)
			if err2 == nil{
				for _,val := range subkey{
					query_reg_value(key,full_path+"\\"+val,value,r)
					}
			}


		}
		
	}else{ //If query is not recursive
		if value!=""{
			test:=display_values(k,value, full_path)
			fmt.Println(test)
		}else{
			value,err1 := k.ReadValueNames(-1)
			if err1 != nil{
				for _,val := range value{
					log.Println(val)
					}
			}
			subkey,err := k.ReadSubKeyNames(-1)
			if err != nil{
				for _,val := range subkey{
					fmt.Println(full_path+"\\"+val)
				}
			}
			
			
		}

	}
	return 0
}

//Function to add subkeys and or values
func add_reg_value (key registry.Key,full_path string,value string, data_type string,data string)(custom_error int){



	first := strings.Split(full_path,"\\")
	var relative_path string=(full_path[len(first[0])+1:])
	k,_,err := registry.CreateKey(key, relative_path, registry.CREATE_SUB_KEY|registry.SET_VALUE)
	if err != nil {
		//error types definition
		fmt.Println("cannot create registry subkey")
		return 1
	}
	defer k.Close()	

	switch strings.ToUpper(data_type){
	case "REG_QWORD":
		input,err:=strconv.ParseUint(data,10,32)
		if err != nil{
			fmt.Println("Please enter Integer value")
			return 1
		}
		err2 := k.SetQWordValue(value,input)
		if err2 != nil{
			log.Println(err)
			return 1
		}
	case "REG_EXPAND_SZ": 
		err2 := k.SetExpandStringValue(value,data)
		if err2 != nil{
			log.Println(err)
			return 1
		}
	case "REG_SZ": 
		err2 := k.SetStringValue(value,data)
		if err2 != nil{
			log.Println(err)
			return 1
		}
	case "REG_BINARY":
		hex,err:=hex.DecodeString(data)
		if err != nil{
			fmt.Println("Please enter hex value")
			return 1
		}
		err2 := k.SetBinaryValue(value,hex)
		if err2 != nil{
			log.Println(err)
			return 1
		}
	case "REG_DWORD":
		input,err:=strconv.ParseUint(data,10,32)
		if err != nil{
			fmt.Println("Please enter Integer value")
			return 1
		}
		err2 := k.SetDWordValue(value,uint32(input))
		if err2 != nil{
			log.Println(err)
			return 1
		}
	case "REG_MULTI_SZ":
		input:=strings.Split(data,"\\0")
		err2 := k.SetStringsValue(value,input)
		if err2 != nil{
			log.Println(err)
			return 1
		}
	case "REG_NONE": 
		return
	default:
		fmt.Printf("SPECIFIED DATA TYPE IS NOT SUPPORTED  \n")
		fmt.Println("RegKey data types \n [ REG_SZ    | REG_MULTI_SZ | REG_EXPAND_SZ | \n REG_DWORD | REG_QWORD    | REG_BINARY    | REG_NONE ] \n If omitted, REG_SZ is assumed.")
		return 1

	}

	return 0
}

func delete_reg_key(key registry.Key,path string)(custom_error int){

	first := strings.Split(path,"\\")
	var relative_path string=(path[len(first[0])+1:])
	err:=registry.DeleteKey(key,relative_path)
	if err!= nil{
		log.Println(err)
		return 1
	}
	fmt.Println(path+" was sucessfully deleted")
	return 0
}
func main(){

key_dict := map[string]registry.Key {"HKLM":registry.LOCAL_MACHINE,"HKCU":registry.CURRENT_USER,"HKCR":registry.CLASSES_ROOT,"HKU":registry.USERS,"HKCC":registry.CURRENT_CONFIG}

if len(os.Args) < 2 {
	fmt.Printf ("no argumets supplied \n")
	fmt.Printf("use -h for help")
	os.Exit(0)
}

add_reg:=flag.String("add","","Enter the path with the add flag to add to registry.\nEnter the value type and data of the registry.\nEx: reg -add HKLM\\Software\\MyCo -v Data -t REG_BINARY -d fe340ead ")
del_reg:=flag.String("delete","","Enter the path with the delete flag to delete from registry.\nEx: reg -delete HKLM\\Software\\MyCo ")
query_reg:=flag.String("query","","Enter the path with the query flag to query from registry\nEx: reg -query HKLM\\Software\\MyCo \n use -s flag to search recursively \n You can search for certain values usinf -v flag")

value := flag.String("v","","Enter the  value of the registry")
data_type := flag.String("t","REG_SZ","Enter the type of the data to be stored")
data := flag.String("d","00","Enter the data to be inserted")
recursive_check := flag.Bool("s",false,"Perform recusive check to query the registry for subkeys and values")

flag.Parse()


var err int 
var full_path string
if *add_reg!=""{
	full_path=*add_reg
	err=add_reg_value(key_dict[strings.ToUpper(strings.Split(full_path,"\\")[0])],full_path,*value,*data_type,*data,)
} else if *del_reg!=""{
	full_path=*del_reg
	err=delete_reg_key(key_dict[strings.Split(full_path,"\\")[0]],full_path)
} else if *query_reg!=""{
	full_path=*query_reg
	err=query_reg_value(key_dict[strings.Split(full_path,"\\")[0]],full_path,*value,*recursive_check)

} else {
	fmt.Println("Please see the help menu, you can add, delete or query the keys")
	os.Exit(0)
}

if(err==1){
	fmt.Println("\nTask Unsucessfull")
}else{
	fmt.Println("\nTask Sucessfull")
}
	
}