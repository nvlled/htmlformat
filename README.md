# htmlformat

A (WIP) go library for formatting HTML.

## Restrictions
- (super)long lines are never broken, and left as it is
- attributes are always on a single line after the tag name

## Why

There are already plenty of existing HTML formatters.
The notable difference in this one is that
the surrounding whitespaces between and around the nodes 
are maintained.

For instance:
```
<p>
        Whitespaces around, between, in, after texts/nodes
should be <i>maintained</i>. The period from
    the last sentence should not be separated with
    a space, but this one <b>is</b> .
    <em> This has a leading space</em>
        <em>This has a trailing space </em>
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
</p>
```

While the HTML may look organized, whitespaces
are inserted where it shouldn't be.
This will look weird when rendered,
most visible is the period separated by a space.

Whereas, this library will output:
```
<p>
    Whitespaces around, between, in, after texts/nodes
    should be <i>maintained</i>. The period from
    the last sentence should not be separated with
    a space, but this one <b>is</b> .
    <em> This has a leading space</em>
    <em>This has a trailing space </em>
</p>
```
