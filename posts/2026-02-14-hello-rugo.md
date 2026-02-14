# Hello, Rugo! 💖

Welcome to the very first post on this blog — and what better way to kick things off than with a love letter to a little language that stole my heart?

Rugo is a Ruby-inspired language that compiles to native binaries through Go. It brings the warmth and expressiveness of Ruby syntax to the world of fast, standalone executables. No runtime to install, no dependencies to chase — just pure joy, compiled.

## Getting Started

Rugo requires Go 1.24+. Install it with a single command:

<img src="terminal-install.svg" alt="Installing Rugo from a terminal" style="max-width:100%;border-radius:10px;margin:1.5rem 0;">

## Your First Rugo Program

Every great journey starts with a hello:

```ruby
puts "Hello, Rugo! 💖"
```

```text
Hello, Rugo! 💖
```

Simple, clean, familiar. If you've ever written Ruby, you'll feel right at home.

## A Language That Feels Like a Hug

Rugo keeps things cozy. Variables, strings, and loops all work the way you'd expect:

```ruby
languages = ["Ruby", "Go", "Rugo"]
for i, lang in languages
  if lang == "Rugo"
    puts "#{lang} — the best of both worlds!"
  else
    puts "#{lang} is lovely too."
  end
end
```

```text
Ruby is lovely too.
Go is lovely too.
Rugo — the best of both worlds!
```

## Functions with Heart

Defining functions is as sweet as it gets:

```ruby
def love(thing)
  return "I love #{thing}!"
end

feelings = ["simplicity", "native binaries", "Ruby syntax"]
for feeling in feelings
  puts love(feeling)
end
```

```text
I love simplicity!
I love native binaries!
I love Ruby syntax!
```

## Hashes and Sweet Data

Rugo handles structured data with grace. Colon syntax creates clean hashes with dot access — no brackets needed:

```ruby
valentines = [
  {name: "Expressiveness", hearts: 5},
  {name: "Speed", hearts: 4},
  {name: "Simplicity", hearts: 5}
]

for i, v in valentines
  hearts = ""
  for j in [1, 2, 3, 4, 5]
    if j <= v.hearts
      hearts += "💖"
    end
  end
  puts "#{v.name}: #{hearts}"
end
```

```text
Expressiveness: 💖💖💖💖💖
Speed: 💖💖💖💖
Simplicity: 💖💖💖💖💖
```

## A Taste of Gummy

Rugo has a growing ecosystem of libraries. [Gummy](https://github.com/rubiojr/gummy) is a tiny ORM that wraps SQLite with a clean, expressive API — define models, insert records, and query with no SQL strings:

```ruby
require "github.com/rubiojr/gummy" with db

conn = db.open(":memory:")
Users = conn.model("users", {name: "text", age: "integer"})

alice = Users.insert({name: "Alice", age: 30})
puts "#{alice.name} is #{alice.age}"

bob = Users.insert({name: "Bob", age: 25})
count = Users.tally(nil)
puts "#{count} users in the database"

conn.close()
```

```text
Alice is 30
2 users in the database
```

Models return smart records with dot access — no boilerplate, just data.

Gummy also has built-in full text search powered by SQLite FTS5:

```ruby
require "github.com/rubiojr/gummy" with db

conn = db.open(":memory:")
Posts = conn.model("posts", {title: "text", body: "text"})
Posts.searchable(["title", "body"])

Posts.insert({title: "Hello Rugo", body: "A lovely intro to the language"})
Posts.insert({title: "CLI Tools", body: "Building command-line apps"})
Posts.insert({title: "Gummy ORM", body: "A tiny ORM for Rugo with SQLite"})

results = Posts.search("rugo", nil)
for post in results
  puts post.title
end
conn.close()
```

```text
Hello Rugo
Gummy ORM
```

## Compile and Share the Love

The magic of Rugo is that your scripts become standalone native binaries:

```bash
rugo build hello.rugo -o hello
./hello  # runs anywhere — no runtime needed!
```

That's it. Write something beautiful, compile it, and share it with the world.

## Why Rugo?

Rugo sits in a sweet spot: the elegance and readability of Ruby, the performance and portability of compiled Go binaries. Whether you're writing a quick script or building a CLI tool, Rugo makes it a delight.

This is just the beginning. Stay tuned for more posts exploring what Rugo can do. Until then — happy coding, and happy Valentine's Day! 💖
