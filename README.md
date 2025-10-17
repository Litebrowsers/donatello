# Donatello
This system is inspired by ideas from [Google’s Picasso](https://dl.acm.org/doi/pdf/10.1145/2994459.2994467) research 
and uses low-level fingerprinting techniques.

## Problem Description
When a user requests a test task from the server, they receive two tasks:

1. A task for rendering a **stable canvas** (baseline test).
2. A task for rendering a **canvas for subpixel analysis**.

On the client side, the user evaluates these tasks and sends the results back to the
server.

## Server-Side Logic

## Canvas Task Encoding Format

This format is used to describe shapes that should be rendered on a canvas.  
Each shape is encoded as a string, and multiple shapes can be combined using `;` as a separator.

The first token is the **shape type**, followed by parameters specific to that shape.  
All colors are specified in **hexadecimal RGB** (`RRGGBB`).

## Shape Types

### 1. Rectangle
```

R:COLOR:W:H:X:Y

```
- `COLOR` – fill color (hex)
- `W` – width
- `H` – height
- `(X,Y)` – top-left corner position

**Example:**
```

R:FF0000:5:3:10:5

```
→ Red rectangle, 5×3, positioned at (10,5).

---

### 2. Circle
```

C:COLOR:R:X:Y

```
- `COLOR` – fill color (hex)
- `R` – radius
- `(X,Y)` – center position

**Example:**
```

C:00FF00:4:15:15

```
→ Green circle, radius 4, centered at (15,15).

---

### 3. Triangle
```

T:COLOR:X1:Y1:X2:Y2:X3:Y3

```
- `COLOR` – fill color (hex)
- `(X1,Y1), (X2,Y2), (X3,Y3)` – vertices of the triangle

**Example:**
```

T:0000FF:2:2:6:2:4:6

```
→ Blue triangle with vertices at (2,2), (6,2), (4,6).

---

### 4. Line
```

L:COLOR:X1:Y1:X2:Y2:THICKNESS

```
- `COLOR` – stroke color (hex)
- `(X1,Y1), (X2,Y2)` – start and end points
- `Thickness` - thickness of line

**Example:**
```

L:FF00FF:5:5:12:8:2

```
→ Purple line from (5,5) to (12,8) with thickness equal to 2.

---

### 5. Ellipse
```

E:COLOR:RX:RY:X:Y

```
- `COLOR` – fill color (hex)
- `RX`, `RY` – radii along X and Y axes
- `(X,Y)` – center position

**Example:**
```

E:FFFF00:6:3:20:10

```
→ Yellow ellipse with radii 6×3, centered at (20,10).

---

## Combining Shapes
Multiple shapes can be combined into a single task string, separated by `;`.

**Example:**
```

R:FF0000:5:3:10:5;C:00FF00:4:15:15;T:0000FF:2:2:6:2:4:6

```

## Storage
The current implementation uses GORM for database interactions and an in-memory cache for temporary data storage.
The system is designed to be extensible, allowing for the future addition of other storage backends such as Redis for 
caching or PostgreSQL for the main database.

In the current approach, a single test is generated and sent to the client as a task. This test consists of a set of 
randomly generated shapes that the client must render. The client then calculates a hash of the rendered output and 
sends it back to the server for verification. This method allows for a baseline analysis of the client's rendering 
capabilities.