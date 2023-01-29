# AWS Report (for specific purposes)

用於協助 DevOps 工作上的日常報表作業。

## Installing

Build the binary file.
```
git clone git@github.com:rickchangch/aws-report.git
cd aws-report
go build -o aws-report .

./aws-report -h
```

Or install it at your go/bin path.
```
go install github.com/rickchangch/aws-report

aws-report -h
~/go/bin/aws-report -h
```

## Usage

Weekly report
```
aws-report weekly -a true -d x.x.x.x -s 2023-01-01 -e 2023-01-14
```

Monthly report
```
aws-report monthly -f "../costs.csv" -d 6
```
