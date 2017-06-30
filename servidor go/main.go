package main

import (
	//conexion a la base de datos
	"database/sql"
	//conversiones de String a otro tipo
	"strconv"
	//manejo de logs
	"log"
	
	//framework web
	"github.com/gin-gonic/gin"
	//genera automaticamente tablas en la BD
	"github.com/coopernurse/gorp"	
	//Conexion a base de datos mysql
	_ "github.com/go-sql-driver/mysql"
)

const (
	//host
	DB_HOST = "localhost"
	//nombre de la BD
	DB_NAME = "autos"
	//Usuario de la BD
	DB_vehiculo = "postgres"
	//Contraseña de la BD
	DB_PASS = "postgres"
)

//Base de datos del usuario
type Vehiculo struct {
	//Primary key de la tabla
	Id        int64  `db:"id" json:"id"`
	//Placa del vehiculo
	Placa string `db:"Placa" json:"placa"`
	//Marca del vehiculo
	Marca  string `db:"Marca" json:"marca"`
	//Linea del vehiculo
	Linea  string `db:"Marca" json:"linea"`
	//Modelo del vehiculo
	Modelo  string `db:"Marca" json:"modelo"`
}


//Relacion entre los controladores
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(Cors())
	//GENERA RUTAS DE PETICIONES
	v1 := r.Group("api/v1")
	{
		v1.GET("/vehiculos", Getvehiculos)
		v1.GET("/vehiculos/:id", Getvehiculo)
		v1.POST("/vehiculos", Postvehiculo)
		v1.PUT("/vehiculos/:id", Updatevehiculo)
		v1.DELETE("/vehiculos/:id", Deletevehiculo)
		v1.OPTIONS("/vehiculos", Optionsvehiculo)     // POST
		v1.OPTIONS("/vehiculos/:id", Optionsvehiculo) // PUT, DELETE
	}
	r.Run(":8080")
}

//Representa la BD
var dbmap = initDb()

//CONECTA A LA BD
func initDb() *gorp.DbMap {
	dsn := DB_vehiculo + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(vehiculo{}, "vehiculo").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

//Verifica errores
func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

//trae los vehiculos de la tabla vehiculo
func Getvehiculos(c *gin.Context) {
	var vehiculos []vehiculo
	_, err := dbmap.Select(&vehiculos, "SELECT * FROM vehiculo")
	if err == nil {
		c.JSON(200, vehiculos)
	} else {
		c.JSON(404, gin.H{"error": "no vehiculo(s) dentro de la tabla"})
	}
	// curl -i http://localhost:8080/api/v1/vehiculos
}

//trae un vehiculo de la tabla vehiculo
func Getvehiculo(c *gin.Context) {
	id := c.Params.ByName("id")
	var vehiculo Vehiculo
	err := dbmap.SelectOne(&vehiculo, "SELECT * FROM vehiculo WHERE id=?", id)
	if err == nil {
		veh_id, _ := strconv.ParseInt(id, 0, 64)
		content := &Vehiculo{
			Id:        veh_id,
			Placa: vehiculo.Placa,
			Marca:  vehiculo.Marca,
			Linea: vehiculo.Linea,
			Modelo: vehiculo.Modelo,

		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "vehiculo no encontrado"})
	}
	// curl -i http://localhost:8080/api/v1/vehiculos/1
}


//Agrega un vehiculo a la tabla vehiculo
func InsertarVehiculo(c *gin.Context) {
	var vehiculo Vehiculo
	c.Bind(&vehiculo)
	if vehiculo.Placa != "" && vehiculo.Marca != "" && vehiculo.Linea != "" && vehiculo.Modelo != ""{
		if insert, _ := dbmap.Exec(`INSERT INTO vehiculo (placa,marca,linea,modelo) VALUES (?, ?, ?, ?)`, vehiculo.Placa, vehiculo.Marca, vehiculo.Linea, vehiculo.Modelo); insert != nil {
			veh_id, err := insert.LastInsertId()
			if err == nil {
				content := &Vehiculo{
					Id:        veh_id,
					Placa: vehiculo.Placa,
					Marca:  vehiculo.Marca,
					Linea: vehiculo.Linea,
					Modelo: vehiculo.Modelo,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}
	} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"Placa\": \"Thea\", \"Marca\": \"Queen\" }" http://localhost:8080/api/v1/vehiculos
}

//Actualiza los datos de un vehiculo de la tabla vehiculo
func Updatevehiculo(c *gin.Context) {
	id := c.Params.ByName("id")
	var vehiculo vehiculo
	err := dbmap.SelectOne(&vehiculo, "SELECT * FROM vehiculo WHERE id=?", id)
	if err == nil {
		var json vehiculo
		c.Bind(&json)
		vehiculo_id, _ := strconv.ParseInt(id, 0, 64)
		vehiculo := vehiculo{
			Id:        vehiculo_id,
			Placa: json.Placa,
			Marca:  json.Marca,
			Linea: json.Linea, 
			Modelo: json.Modelo,
		}
		if vehiculo.Placa != "" && vehiculo.Marca != "" {
			_, err = dbmap.Update(&vehiculo)
			if err == nil {
				c.JSON(200, vehiculo)
			} else {
				checkErr(err, "Updated failed")
			}
		} else {
			c.JSON(422, gin.H{"error": "fields are empty"})
		}
	} else {
		c.JSON(404, gin.H{"error": "vehiculo not found"})
	}
	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"Placa\": \"Thea\", \"Marca\": \"Merlyn\" }" http://localhost:8080/api/v1/vehiculos/1
}

//Elime¿ina un registro de la tabla vehiculo
func Deletevehiculo(c *gin.Context) {
	id := c.Params.ByName("id")
	var vehiculo vehiculo
	err := dbmap.SelectOne(&vehiculo, "SELECT id FROM vehiculo WHERE id=?", id)
	if err == nil {
		_, err = dbmap.Delete(&vehiculo)
		if err == nil {
			c.JSON(200, gin.H{"id #" + id: " deleted"})
		} else {
			checkErr(err, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "vehiculo not found"})
	}
	// curl -i -X DELETE http://localhost:8080/api/v1/vehiculos/1
}


//Garnatiza las opciones insertar actualizar y borrar
func Optionsvehiculo(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
