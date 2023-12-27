# Sparse Set XCC solver in Go
[Sparse-set XCC Solver](https://www-cs-faculty.stanford.edu/~knuth/programs/ssxcc.w) in Go based on Donald Knuth's Dancing Cells presented in the [27th Annual Christmas Tree Lecture](https://www.youtube.com/watch?v=622iPkJfYrI&list=WL&index=3).

## Problem statement
The solver takes:

* a set of primary items;
* a set of secondary items;
* a set of options, which are subsets of the primary and secondary items.

The solver’s job is to find a subset of the options that

* includes each primary item once and only once, and
* colors each secondary item consistently.
Options can contain secondary items with or without colors. If a secondary item has no color, then the solver will not use it more than once (so that it defines a “zero or one” constraint). If an option has a secondary item with a color, then the solver can use it with the same color as many times as it wants, but not uncolored or with a different color.

### Toy Example
```
p q r
x y
  
p q x y:A
p r x:A y
p x:B
q x:A
r y:B
```

## Exact Covering with Colors (XCC)
It's an extension of the "exact cover" problem -- you're given a set of items that should be covered, and a set of options (sets of items) that may cover some items. Items are divided into two groups:
* Primary: these items need to be covered exactly once.
* Secondary: these items may be either uncovered, or have possibly-many covering options, but the "color" of the coverings must be the same.

[Here](https://docs.rs/xcc/latest/xcc/) you can find a more detailed explanation of XCC problems, however, they do not use the *"exact covering with colors"* used by Knuth. 
> ToDo: find a better reference, perhaps a link to The Art of Computer Programming book.

## Sparse Set
A sparse set is a simple data structure that has following properties:
* O(1) to add an item.
* O(1) to remove an item.
* O(1) to lookup an item.
* O(1) to clear the set.
* O(n) to iterate over the set.
* Does not require initializing its internal items storage.

Solving this problem on its own is an interesting programming exercise. For more information refer to these two blog posts:
* [Sparse sets by Oleksandr Manenko](https://manenko.com/2021/05/23/sparse-sets.html)
* [Using Uninitialized Memory for Fun and Profit by Russ Cox](https://research.swtch.com/sparse)
 


