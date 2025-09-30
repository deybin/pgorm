# Proyecto PGORM

PGORM es una librería en Go para construir y ejecutar consultas SQL de forma dinámica utilizando [pgx](https://github.com/jackc/pgx).

## Descripción

PGORM facilita la interacción con bases de datos PostgreSQL a través de `pgx`. Permite construir consultas SQL mediante métodos encadenados (fluent API) para establecer diferentes partes de la consulta como `SELECT`, `WHERE`, `JOIN`, `ORDER BY`, etc.  
Además, proporciona validaciones de datos antes de ejecutar operaciones de inserción, actualización o eliminación, y soporta transacciones para garantizar la consistencia.

## Características

- **Ejecución de consultas SQL**: simple y segura sobre PostgreSQL.
- **Soporte para cláusulas comunes**: `SELECT`, `WHERE`, `JOIN`, `ORDER BY`, etc.
- **Consultas personalizadas**: flexibilidad para ejecutar SQL directo.
- **Manejo de errores integrado**: control de errores SQL en cada paso.
- **Validación de datos**: antes de `INSERT`, `UPDATE` o `DELETE` según reglas del esquema (longitud, tipo de dato, valores permitidos, etc.).
- **Transacciones**: soporte para agrupar múltiples operaciones en una misma unidad de trabajo.

## Instalación

Para instalar PGORM, simplemente ejecuta:

````bash
go get github.com/deybin/pgorm

````

## Uso

- variables de entorno

  - ENV_DDBB_SERVER='IP/Host'
  - ENV_DDBB_USER='user'
  - ENV_DDBB_PASSWORD='password'
  - ENV_DDBB_DATABASE='database_default'
  - ENV_DDBB_PORT=port
  - ENV_KEY_CRYPTO='key_para_encriptar'

```go
package main

import (
	"fmt"
	"github.com/deybin/pgorm"
)

func main() {
	result, err := new(pgorm.Query).New(pgorm.QConfig{Database: "mi_database"}).SetTable("mi_tabla").
		Select().Where("campo", pgorm.I, "valor_campo").And("campo2", pgorm.IN, []interface{}{"valor1", "valor2", "valor3"}).OrderBy("campo_ordenar DESC").
		Exec().All()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Resultado:", result)
}
```

## Ejemplo con Transacciones

PGORM también soporta transacciones de manera sencilla.  
Puedes agrupar varias operaciones (`INSERT`, `UPDATE`, `DELETE`) en una misma unidad de trabajo y confirmar los cambios con `Commit()` o revertirlos con `Rollback()`.

```go
package main

import (
	"fmt"
	"github.com/deybin/pgorm"
)

func main() {
	// Iniciamos la transacción
	tx := pgorm.SqlExecMultiple{}
	tx.New("mi_database")

	// Primer INSERT
	dataInsert :=  map[string]interface{}{
		"nombre": "Juan Pérez",
		"email":  "juan@example.com",
	}

	// Segundo INSERT
		dataInsert2:= map[string]interface{}{
		"cliente_id": 1,
		"saldo":      5000,
	}

tx.SetInfo(new(tables.Schema_Model1).New(), dataInsert)
tx.SetInfo(new(tables.Schema_Model2).New(), dataInsert2)

	if err := crud.Exec(); err != nil {
		fmt.Println("Error:", err)
	}
}
```

## Contribución

¡Las contribuciones son bienvenidas! Si quieres contribuir a este proyecto o encuentras algún problema por favor abre un issue primero para discutir los cambios propuestos.

## Licencia

© Deybin, 2025

Este proyecto está bajo la Licencia MIT. Consulta el archivo [LICENSE](https://github.com/deybin/pgorm/blob/master/LICENCE) para más detalles.
