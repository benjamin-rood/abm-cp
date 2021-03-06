function DrawList(obj) {
  this.cpPrey = obj.data.cpPrey
  this.vp = obj.data.vp
  this.bg = obj.data.bg
}

var initDrawObj = {
  type: 'drawlist',
  data: {
    cpPrey: [],
    vp: [],
    bg: {
      red: 0,
      green: 0,
      blue: 0
    }
  }
}

var sketch = function(p) {
  var modelDw = 1.0
  var modelDh = 1.0
  p.setup = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.75
    p.createCanvas(w, h)
    p.noLoop()
    p.background(0, 0, 255)
  }

  p.draw = function() {
    p.background(drawlist.bg.red, drawlist.bg.green, drawlist.bg.blue)

    if drawlist.cpPrey {
      for (var i = 0; i < drawlist.cpPrey.length; i++) {
        var x = absToView(drawlist.cpPrey[i].position.x, modelDw, p.width)
        var y = absToView(drawlist.cpPrey[i].position.y, modelDh, p.height)
        var col = p.color(drawlist.cpPrey[i].colour.red, drawlist.cpPrey[i].colour.green, drawlist.cpPrey[i].colour.blue)
        p.strokeWeight(8)
        p.stroke(col)
        p.point(x, y)
      }
    }

    if drawlist.vp {
      for (var i = 0; i < drawlist.vp.length; i++) {
        var x = drawlist.vp[i].position.x
        var y = drawlist.vp[i].position.y
        var angle = 0 //  because positive rotations here happen clockwise, rather than the convention of the unit circle
        var col = p.color(drawlist.vp[i].colour.red, drawlist.vp[i].colour.green, drawlist.vp[i].colour.blue)
        p.push()
        p.translate(x, y)
        p.rotate(angle)
        p.fill(col)
        p.strokeWeight(1)
        p.stroke(255)
        p.ellipse(0, 0, vSize, vSize)
          // p.triangle(0-tSize, 0+tSize, 0, 0-(2*tSize), 0+tSize, 0+tSize)
        p.pop()
      }
    }
  }

  p.windowResized = function() {
    var w = $('#abm-viewport').innerWidth()
    var h = w * 0.75
    p.resizeCanvas(w, h)
  }
}

function absToView(p, d, n) {
  view = (((p + d) / (2 * d)) * n)
  return view
}
