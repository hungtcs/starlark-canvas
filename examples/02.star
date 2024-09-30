S = 1024
dc = canvas.Context(S, S)
dc.set_rgba(1, 0, 0, 0.1)
for i in range(0, 360, 15):
    dc.push()
    dc.rotate_about(canvas.radians(i), S/2, S/2)
    dc.draw_ellipse(S/2, S/2, S*7/16, S/8)
    dc.fill()
    dc.pop()

dc.set_fill_rule(0)
dc.set_fill_rule(1)

print(dc.get_data_uri())
