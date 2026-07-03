package sample

import "github.com/lcnascimento/go-kit/env"

func Getter() {
	// envvar.Get busca uma variável de ambiente e converte para o tipo desejado
	name := env.Get[string]("NAME")
	println(name)

	// envvar.Get suporta um valor padrão com a opção WithDefaultValue
	age := env.Get[int]("AGE", env.WithDefaultValue(30))
	println(age)
}
