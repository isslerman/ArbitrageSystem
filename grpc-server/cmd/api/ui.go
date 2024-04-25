package main

// //	START UI TERM
// //
// // //////////
// // type Panel struct {
// // }

// func (app *Config) initUI() {
// 	if err := ui.Init(); err != nil {
// 		log.Fatalf("failed to initialize termui: %v", err)
// 	}
// 	defer ui.Close()

// 	app.p = widgets.NewParagraph()
// 	app.p.Title = "Text Box"
// 	// app.p.Text = "PRESS q TO QUIT DEMO"
// 	app.p.SetRect(0, 0, 50, 5)
// 	app.p.TextStyle.Fg = ui.ColorWhite
// 	app.p.BorderStyle.Fg = ui.ColorCyan

// 	// updateParagraph := func(count int) {
// 	// 	if count%2 == 0 {
// 	// 		p.TextStyle.Fg = ui.ColorRed
// 	// 	} else {
// 	// 		p.TextStyle.Fg = ui.ColorWhite
// 	// 	}
// 	// }

// 	// isrunning := true

// 	// // for isrunning {
// 	// // 	// uiEvents := ui.PollEvents()
// 	// // 	ui.Render(app.p)
// 	// // 	time.Sleep(time.Millisecond * 20)
// 	// // }

// 	uiEvents := ui.PollEvents()
// 	ticker := time.NewTicker(time.Second).C
// 	tickerCount := 1
// 	for {
// 		select {
// 		case e := <-uiEvents:
// 			switch e.ID {
// 			case "q", "<C-c>":
// 				return
// 			}
// 		case <-ticker:
// 			// updateParagraph(tickerCount)
// 			// draw(tickerCount)
// 			ui.Render(app.p)
// 			tickerCount++
// 		}
// 	}
// 	// tickerCount := 1
// 	// ticker := time.NewTicker(time.Second).C
// 	// for {
// 	// 	select {
// 	// 	case e := <-uiEvents:
// 	// 		switch e.ID {
// 	// 		case "q", "<C-c>":
// 	// 			return
// 	// 		}
// 	// 	case <-ticker:
// 	// 		//updateParagraph(tickerCount)
// 	// 		tickerCount++
// 	// 	}
// 	// }
// 	// if err := ui.Init(); err != nil {
// 	// 	log.Fatalf("failed to initialize termui: %v", err)
// 	// }
// 	// defer ui.Close()

// 	// isrunning := true

// 	// margin := 2
// 	// pheight := 3

// 	// pticker := widgets.NewParagraph()
// 	// pticker.Title = "Binancef"
// 	// pticker.Text = "[BTCUSDT](fg:cyan)"
// 	// pticker.SetRect(0, 0, 14, pheight)

// 	// pprice := widgets.NewParagraph()
// 	// pprice.Title = "Market price"
// 	// ppriceOffset := 14 + 14 + margin + 2
// 	// pprice.SetRect(14+margin, 0, ppriceOffset, pheight)

// 	// pfund := widgets.NewParagraph()
// 	// pfund.Title = "Funding rate"
// 	// pfund.SetRect(ppriceOffset+margin, 0, ppriceOffset+margin+16, 3)

// 	// tob := widgets.NewTable()
// 	// out := make([][]string, 20)
// 	// for i := 0; i < 20; i++ {
// 	// 	out[i] = []string{"n/a", "n/a"}
// 	// }
// 	// tob.TextStyle = ui.NewStyle(ui.ColorWhite)
// 	// tob.SetRect(0, pheight+2, 30, 22+pheight+2)
// 	// tob.PaddingBottom = 0
// 	// tob.PaddingTop = 0
// 	// tob.RowSeparator = false
// 	// tob.TextAlignment = ui.AlignCenter
// 	// for isrunning {
// 	// 	// out[i] = []string{fmt.Sprintf("[%.2f](fg:red)", asks[i].Price), fmt.Sprintf("[%.2f](fg:cyan)", asks[i].Volume)}
// 	// 	// out[i+10] = []string{fmt.Sprintf("[%.2f](fg:green)", bids[i].Price), fmt.Sprintf("[%.2f](fg:cyan)", bids[i].Volume)}

// 	// 	tob.Rows = out

// 	// 	pprice.Text = "10"
// 	// 	pfund.Text = fmt.Sprintf("[%s](fg:yellow)", "fundingRate")
// 	// 	ui.Render(pticker, pprice, pfund, tob)
// 	// 	time.Sleep(time.Millisecond * 20)
// 	// }
// }
