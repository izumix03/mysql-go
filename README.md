# mysql-go
## Installation
Simple install the package to your $GOPATH with the go tool from shell:
```
$ go get -u github.com/go-sql-driver/mysql
```
Make sure Git is installed on your machine and in your system's PATH.

## Usage

```$xslt

type User struct {
	Name       string    `name:"name"`
	Status     int       `name:"status"`
	Created_at time.Time `name:"created_at"`
}

// Table provides this entity's table name
func (e *User) Table() string {
	return `users`
}

var users []User
if err := mysqlgo.WithConnection(
    `root:@tcp(localhost:3306)/company?charset=utf8&parseTime=True`,
    func(con *mysqlgo.Connection) error {
        return con.Where(`status = ?`, 0).
            Where(`created_at < now() - interval 10 minute`).
            OrderBy(`name`).
            FindMany(&users)
    }); err != nil {
    log.Printf(`Failed to get user list. %s`, err.Error())
}

if len(users) > 0 {
    user := users[0]
    log.Println(user.Name) // `taro`
}
```