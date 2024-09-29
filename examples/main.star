
data = [100, 200, 50, 400, 155, 300]
categories = ["原油", "汽油", "柴油", "电力", "氢气", "核能"]

# 数据条数
count = len(data)

def draw():
  # 数据最大值
  maxValue = data[0]
  for item in data:
    if item > maxValue:
      maxValue = item


  # 数据最小值
  minValue = data[0]
  for item in data:
    if item < minValue:
      minValue = item

  width = 1200
  height = 720

  # top right bottom left
  padding = [80, 80, 80, 80]
  paddingTop = padding[0]
  paddingRight = padding[1]
  paddingBottom = padding[2]
  paddingLeft = padding[3]

  contentWidth = width - paddingLeft - paddingRight
  contentHeight = height - paddingTop - paddingBottom

  dc = canvas.Context(width, height)

  # 背景色
  dc.set_hex_color("#FFFFFF")
  dc.draw_rectangle(0, 0, width, height)
  dc.fill()

  # 网格线
  if True:
    # 横向网格数量
    hc = 5
    # 纵向网格数量
    vc = count
    w = contentWidth / vc
    h = contentHeight / hc
    dc.new_sub_path()

    # 横向网格线
    for i in range(0, hc+1):
      dc.move_to(paddingLeft, paddingTop+(h*i))
      dc.line_to(width-paddingRight, paddingTop+(h*i))

    # 纵向网格线
    for i in range(0, vc+1):
      dc.move_to(paddingLeft+w*i, paddingTop)
      dc.line_to(paddingLeft+w*i, height-paddingBottom)

    dc.set_hex_color("#CCCCCC")
    dc.stroke()

  # 坐标轴
  if True:
    dc.new_sub_path()
    dc.set_hex_color("#000000")
    # X 轴
    dc.move_to(paddingLeft, height-paddingBottom)
    dc.line_to(width-paddingRight, height-paddingBottom)
    # Y 轴
    dc.move_to(paddingLeft, height-paddingBottom)
    dc.line_to(paddingLeft, paddingTop)
    dc.stroke()

  # 柱状图
  if True:
    hc = 5
    # 纵向网格数量
    vc = count
    w = contentWidth / vc
    h = contentHeight / hc
    bw = w * 0.618

    # 绘制柱子
    dc.set_hex_color("#00AA00")
    dc.new_sub_path()
    for idx, value in enumerate(data):
      bh = value / maxValue * contentHeight
      dc.draw_rectangle(
        paddingLeft+(w*idx)+w/2-bw/2,
        paddingTop+contentHeight-bh,
        bw,
        bh,
      )
    dc.fill()

    # 绘制标签
    dc.set_hex_color("#000000")
    dc.new_sub_path()
    # dc.set_font_face(defaultFontFace)
    for idx, category in enumerate(categories):
      sw, sh = dc.measure_string(category)
      dc.draw_string(
        category,
        paddingLeft+(w*idx)+w/2-sw/2,
        height-paddingBottom+sh+12,
      )

    # 绘制数值
    dc.set_hex_color("#000000")
    dc.new_sub_path()
    pv = maxValue / 5
    for i in range(0, hc+1):
      s = str(int(pv*i))
      sw, sh = dc.measure_string(s)
      dc.draw_string(
        s,
        paddingLeft-12-sw,
        height-paddingBottom-i*h+sh/2,
      )


  return dc



dc = draw()

uri = dc.get_data_uri()
print(uri)
