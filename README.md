# htmlformat

A (WIP) go library for formatting HTML.

## Behaviour and restrictions
- long lines are never wrapped, and left as it is
- attributes too are always on a single line after the tag name
- *script*, *style*, *code*, *pre* contents are treated as text,
  the indentation is only modified to align with the parent node 

## Why

There are already plenty of existing HTML formatters.
I wrote this one because there is recurring quirk
that I see with HTML formatters, not just in go packages.
Even in JSX formatters too.
The recurring quirk is that the whitespaces 
are not correctly inserted or preserved.

For instance:
```
<p>
        Whitespaces around, between, in, after texts/nodes
should be <i>conserved</i>. The period from
    the last sentence should not be separated with
    a space, but this one <b>is</b> .
    <em> This has a leading space</em>

        <em>This has a trailing space </em>
    
    [<a>link</a>]
    </p>
```

Some formatters will output:

```
<p>
    Whitespaces around, between, in, after texts/nodes
    should be
    <i>
       maintained
    </i>
    . The period from
    the last sentence should not be separated with
    a space, but this one
    <b>
        is
    </b>
    .
    <em>
        This has a leading space
    </em>
    <em>
        This has a trailing space
    </em>
    [
    <a>
        link
    </a>
    ]
</p>
```

While the HTML may look organized, whitespaces
are inserted where it shouldn't be.
This will look weird when rendered, such
as the period dangling on its own, 

Whereas, this library will output:
```
<p>
    Whitespaces around, between, in, after texts/nodes
    should be <i>maintained</i>. The period from
    the last sentence should not be separated with
    a space, but this one <b>is</b> .
    <em> This has a leading space</em>
    <em>This has a trailing space </em>
    [<a>link</a>]
</p>
```

## TODO
- create a site that compares output of different libraries
- upload at pkg.go.dev
- create an executable program