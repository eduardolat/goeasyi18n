# Go Easy i18n

Go Easy i18n allows you to load your translations in different ways:

## Manually

If you have a small set of translations, you can load them manually as shown in the [Basic Usage Example.](../01-basic-usage/main.go)

## From JSON/YAML as Bytes or Strings

You can fetch your JSON/YAML translations from a database or a remote API and load them into the i18n instance as part of the startup of your program, look at the example to see how to do it:

[Bytes/String Example](./01-from-strings/main.go)

## From JSON/YAML files in your server file system

You can load your translations from JSON/YAML files in your file system, look at the example to see how to do it:

[Files Example](./02-from-files/main.go)

## From JSON/YAML files in `fs.FS` (`embed.FS`) file system

If you want to load your translations from JSON/YAML files embedded in your binary or some other implementation of the `fs.FS` interface, look at the example to see how to do it:

[`embed.FS` Example](./03-from-embed-fs/main.go)
