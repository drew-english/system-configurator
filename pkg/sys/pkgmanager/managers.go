package pkgmanager

import (
	"regexp"
	"text/template"
)

var Managers = map[string]PacakgeManager{
	"apk":    apk,
	"apt":    apt,
	"brew":   brew,
	"dnf":    dnf,
	"snap":   snap,
	"pacman": pacman,
}

// Delcaration of the package managers and their commands.
var (
	apk = &basePackageManager{
		BaseCmd:          "apk",
		AddCmd:           cmd("add"),
		RemoveCmd:        cmd("del"),
		ListCmd:          cmd("list --installed"),
		listParsePattern: re(`^([\w-]+)-(\S+-\S+)`),
		versionTmpl:      tpl("{{.Name}}={{.Version}}"),
	}

	apt = &basePackageManager{
		BaseCmd:          "apt",
		AddCmd:           cmd("install", "-y"),
		RemoveCmd:        cmd("remove"),
		ListCmd:          cmd("list", "--installed"),
		listParsePattern: re(`^([\w-]+)\/.*?\s(\S+)`),
		versionTmpl:      tpl("{{.Name}}={{.Version}}"),
	}

	brew = &basePackageManager{
		BaseCmd:          "brew",
		AddCmd:           cmd("install"),
		RemoveCmd:        cmd("remove"),
		ListCmd:          cmd("list", "--versions"),
		listParsePattern: re(`^([\w-]+)\s(\S+)`),
		versionTmpl:      tpl("{{.Name}}"),
	}

	dnf = &basePackageManager{
		BaseCmd:          "dnf",
		AddCmd:           cmd("install", "-y"),
		RemoveCmd:        cmd("erase"),
		ListCmd:          cmd("list", "--installed"),
		listParsePattern: re(`^(\S+)\.\w+\s+(\S+?)-`),
		versionTmpl:      tpl("{{.Name}}-{{.Version}}"),
	}

	snap = &basePackageManager{
		BaseCmd:          "snap",
		AddCmd:           cmd("install", "--classic"),
		RemoveCmd:        cmd("remove"),
		ListCmd:          cmd("list"),
		listParsePattern: re(``), // TODO: Test on ubuntu machine, docker image unable to run snap
		versionTmpl:      tpl("{{.Name}} --channel={{.Version}}"),
	}

	pacman = &basePackageManager{
		BaseCmd:          "pacman",
		AddCmd:           cmd("-S", "--noconfirm"),
		RemoveCmd:        cmd("-Rscn", "--noconfirm"),
		ListCmd:          cmd("-Q"),
		listParsePattern: re(`^([\w-\.]+)\s(\S+)`),
		versionTmpl:      tpl("{{.Name}}={{.Version}}"),
	}
)

func cmd(args ...string) []string {
	return args
}

func re(s string) *regexp.Regexp {
	return regexp.MustCompile(s)
}

func tpl(s string) *template.Template {
	return template.Must(template.New("pkgVersion").Parse(s))
}
