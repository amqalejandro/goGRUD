package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"text/template"
)

var plantilla = template.Must(template.ParseGlob("plantillas/*"))

func conexionBD() (conexion *sql.DB) {
	driver := "mysql"
	usuario := "root"
	clave := "toor"
	bd := "golandcrud"

	conexion, err := sql.Open(driver, usuario+":"+clave+"@tcp(127.0.0.1)/"+bd)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

func main() {

	http.HandleFunc("/", inicio)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/borrar", borrar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/actualizar", actualizar)

	http.ListenAndServe(":8080", nil)
	log.Println("Servidor Corriendo...")
}

func actualizar(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		id := request.FormValue("id")
		nombre := request.FormValue("nombre")
		correo := request.FormValue("correo")

		conexion := conexionBD()
		update, err := conexion.Prepare("UPDATE empleado set nombre=?,correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		update.Exec(nombre, correo, id)
		http.Redirect(writer, request, "/", 301)
	}

}

func editar(writer http.ResponseWriter, request *http.Request) {
	idEmpleado := request.URL.Query().Get("id")
	conexion := conexionBD()
	registro, err := conexion.Query("SELECT * FROM empleado WHERE id=?", idEmpleado)
	empleado := EmpleadoType{}
	if err != nil {
		panic(err.Error())
	}

	for registro.Next() {
		var id int
		var nombre, correro string
		err = registro.Scan(&id, &nombre, &correro)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correro

		plantilla.ExecuteTemplate(writer, "editar", empleado)
	}

}

func borrar(writer http.ResponseWriter, request *http.Request) {
	idEmpleado := request.URL.Query().Get("id")
	conexion := conexionBD()
	borrar, err := conexion.Prepare("DELETE FROM empleado WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrar.Exec(idEmpleado)
	http.Redirect(writer, request, "/", 301)

}

func insertar(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		nombre := request.FormValue("nombre")
		correo := request.FormValue("correo")

		conexion := conexionBD()
		insertar, err := conexion.Prepare("INSERT INTO empleado(nombre,correo) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertar.Exec(nombre, correo)
		http.Redirect(writer, request, "/", 301)
	}
}

func crear(writer http.ResponseWriter, request *http.Request) {
	plantilla.ExecuteTemplate(writer, "crear", nil)
}

type EmpleadoType struct {
	Id     int
	Nombre string
	Correo string
}

func inicio(writer http.ResponseWriter, request *http.Request) {

	conexion := conexionBD()
	registros, err := conexion.Query("SELECT * FROM EMPLEADO")

	if err != nil {
		panic(err.Error())
	}
	empleado := EmpleadoType{}
	empleados := []EmpleadoType{}

	for registros.Next() {
		var idX int
		var nombre, correo string
		err = registros.Scan(&idX, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = idX
		empleado.Nombre = nombre
		empleado.Correo = correo

		empleados = append(empleados, empleado)
	}

	plantilla.ExecuteTemplate(writer, "inicio", empleados)

}
