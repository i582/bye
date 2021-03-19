To get started with bye, let's take a look at the basic commands.

First of all, the `create` command is used to create the skeleton of the project.

```bash (create.sh)
$ bye create <project-name>
```

The command takes one argument, the name of the project. The command will create a folder with the base skeleton of the
project.

Let's go to the folder with the project using the `cd` command:

```
$ cd <project-name>
```

Now, while in the project folder, run the `serve` command.

```shell (serve.sh)
$ bye serve 
```

In the terminal, you should see the process of generating the site, as well as a message that the site has been launched
on the local server.

```shell (output.sh)
Started local server on port 3005.
Go to http://localhost:3005 to view the site.
```

Go to the specified address. The generated sample site should open.

If all is well, then congratulations, now let's take a closer look.

## Folder structure

After executing the `create` command, the basic skeleton of the project was created. Let's take a look at the structure
of the folders.

```
# `bye.yml` is a config file, we will look at it next.
│   bye.yml
│
# `pages` is the folder in which all page sources are stored.
├───pages
# `index.md` at the root, this is the index page source.
│   │   index.md
│   │
# each page is a separate folder.
│   ├───example
# which contains the `index.md` file with the page markup.
│   │       index.md
│   │
│   └───hello-bye
│           index.md
│
# `static` is a folder with styles and page templates.
└───static
    ├───styles
# `style.css` is the main style file. In it, you can change the styles as you want.
    │       style.css
    │
# `templates` are HTML markup of pages with some data that is inserted at the stage of generation. You can change them however you want.
    └───templates
            file.html
            index.html
            page-part.html
            page-preview.html
            page.html
            tag-page-preview.html
            tags.html
```

## Adding pages

Now let's add a new page.

First of all, you need to create a *new folder* in the `pages` folder, and in it create the `index.md` file. Add some
text to it so that the page is not blank.

```sh (new-page.sh)
mkdir "some-folder"
cd "some-folder"
vim index.md
```

Now let's move on to the settings in the config.

We need the `pages` field, which describes all the pages of the site.

Let's add a new entry.

```yaml (bye.yml)
{
	name: "some-folder",
	title: "Some title",
	tags: [
		"some",
		"tag",
	],
},
```

The name **must match** the name of the previously created folder. The title can be anything. Tags can be anything.

The pages will be something like this:

```yaml (bye.yml)
pages: [
    {
      name: "hello-bye",
      title: "Hello Bye",
      tags: [
          "hello",
          "bye",
      ],
    },
    {
      name: "example",
      title: "Example of page",
      tags: [
          "rust",
      ],
    },
    {
      name: "<folder>",
      title: "Some title",
      tags: [
          "some",
          "tag",
      ],
    },
],
```

Now stop the server (press `Ctrl + C`) and run the `serve` command again.

Go to the site and reload it, you should see a new page.

If everything went well, let's take a look at how to delete pages.

## Deleting pages

It's simple, comment out or delete the entry in the config and the page will not be displayed.

> However, if you delete the folder, but *don't delete the entry*, then you will get an error while generating.