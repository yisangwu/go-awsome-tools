

命令行参数格式：
-db_host aaa / --db_host aaa / -db_host=aaa / --database=aaaa



直接运行代码：

分析文件导入DB并生成json文件：
go run .\QuizeMasterImport.go -db_host="localhost" -db_user=admin -db_password=admin -database=quiz_master -table=question_bank -file_name="quize.txt"


仅仅从DB中生成json文件：
go run .\QuizeMasterImport.go -db_host="192.168.1.222" -db_user=admin -db_password=admin -database=quiz_master -table=question_bank -genx=true

注意:
因直接从excel中导出格式CSV，因带有特殊字符。
先另存为制表符分隔的txt文件，再在记事本中，另存为utf-8格式。