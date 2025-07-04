// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/rvflash/elapsed"
import "strconv"
import "domofonEmulator/client/models"

func SingleDomofonCard(props models.Intercom) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = SingleDomofonCardStyle().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		date := elapsed.LocalTime(props.CreatedAt, "ru")
		doorStatusRu := ""
		if props.DoorStatus == false {
			doorStatusRu = "закрыта"
		} else {
			doorStatusRu = "открыта"
		}
		intercomeStatusRu := ""
		if props.IntercomStatus == false {
			intercomeStatusRu = "Выключен"
		} else {
			intercomeStatusRu = "Включен"
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"domofon-card-container\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 = []any{
			templ.KV("domofon-card active", props.IntercomStatus == true),
			templ.KV("domofon-card inactive", props.IntercomStatus != true)}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var2...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<div class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var2).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\"><div class=\"card-header\"><span class=\"domofon-id\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(props.ID)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 28, Col: 55}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</span> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 = []any{"status-badge"}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var5...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<span class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var5).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(intercomeStatusRu)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 30, Col: 47}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</span><div class=\"bell-and-door--wrapper\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = phoneIcon(props.IsCalling).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "<div class=\"door-status--wrapper\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = doorIcon(props.DoorStatus).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 = []any{
			templ.KV("door-status open", props.DoorStatus == true),
			templ.KV("door-status closed", props.DoorStatus != true)}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var8...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "<span class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 string
		templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var8).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(doorStatusRu)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 39, Col: 46}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</span></div></div></div><div class=\"card-body\"><div class=\"mac-address\"><span class=\"label\">MAC:</span> <span class=\"value\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var11 string
		templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(props.MAC)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 48, Col: 59}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</span></div><div class=\"address\"><span class=\"label\">Адрес:</span> <span class=\"value\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var12 string
		templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(truncateText(props.Address, 50))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 52, Col: 81}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "</span></div></div><div class=\"card-footer\"><div class=\"apartments\"><span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var13 string
		templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(props.NumberOfApartments)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 58, Col: 56}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, " кв.</span></div><div class=\"created-at\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var14 string
		templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(date)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 61, Col: 30}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "</div></div></div><div class=\"domofon-controls\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = ControlButtonStyle().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "<form action=\"/domofon/call\" method=\"post\" class=\"control-form\"><input type=\"hidden\" name=\"domofon_id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var15 string
		templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 68, Col: 89}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "\"><div class=\"form-group\"><label for=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var16 string
		templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs("apartment_" + strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 70, Col: 71}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "\">Квартира:</label> <input type=\"number\" id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var17 string
		templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs("apartment_" + strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 71, Col: 84}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "\" name=\"apartment\" min=\"1\" max=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var18 string
		templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(
			strconv.Itoa(props.NumberOfApartments))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 72, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "\" required class=\"apartment-input\"></div><button type=\"submit\" class=\"control-button call-button\"><svg width=\"16\" height=\"16\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\"><path d=\"M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z\"></path></svg> Позвонить</button></form><form action=\"/domofon/open\" method=\"post\" class=\"control-form\"><input type=\"hidden\" name=\"domofon_id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var19 string
		templ_7745c5c3_Var19, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 83, Col: 89}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var19))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "\"><div class=\"form-group\"><label for=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var20 string
		templ_7745c5c3_Var20, templ_7745c5c3_Err = templ.JoinStringErrs("door_apartment_" + strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 85, Col: 76}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var20))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "\">Квартира:</label> <input type=\"number\" id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var21 string
		templ_7745c5c3_Var21, templ_7745c5c3_Err = templ.JoinStringErrs("door_apartment_" + strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 86, Col: 89}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var21))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "\" name=\"apartment\" min=\"1\" max=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var22 string
		templ_7745c5c3_Var22, templ_7745c5c3_Err = templ.JoinStringErrs(
			strconv.Itoa(props.NumberOfApartments))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 87, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var22))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "\" required class=\"apartment-input\"></div><button type=\"submit\" class=\"control-button open-button\"><img width=\"16px\" height=\"16px\" class=\"dooricon\" src=\"/client/web/static/icons/door.svg\" alt=\"doorinactive ico\"> Открыть дверь</button></form><div class=\"power-controls\"><form hx-encoding=\"multipart/form-data\" hx-post=\"/api/powerIntercom\" hx-target=\"#notification-area\" hx-target-error=\"#notification-area\" hx-trigger=\"submit\" hx-swap=\"innerHTML\" enctype=\"multipart/form-data\" hx-on::after-request=\"if(event.detail.successful) this.reset()\" class=\"control-form\"><input type=\"hidden\" name=\"domofon_id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var23 string
		templ_7745c5c3_Var23, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 99, Col: 93}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var23))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, "\"> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if props.IntercomStatus == true {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, "<input type=\"hidden\" name=\"action\" value=\"off\"> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if props.IntercomStatus == false {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 27, "<input type=\"hidden\" name=\"action\" value=\"on\"> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		var templ_7745c5c3_Var24 = []any{templ.KV("control-button power-buttonON", props.IntercomStatus == true),
			templ.KV("control-button power-buttonOFF", props.IntercomStatus != true)}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var24...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 28, "<button type=\"submit\" class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var25 string
		templ_7745c5c3_Var25, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var24).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var25))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 29, "\"><svg width=\"16\" height=\"16\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\"><path d=\"M18.36 6.64a9 9 0 1 1-12.73 0\"></path> <line x1=\"12\" y1=\"2\" x2=\"12\" y2=\"12\"></line></svg> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if props.IntercomStatus == true {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 30, "<span>Выключить</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 31, "<span>Включить\u2063\u2063\u2063 \u2063\u2063\u2063\u2063</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 32, "</button></form><form action=\"/domofon/delete\" method=\"post\" class=\"control-form delete-form\"><input type=\"hidden\" name=\"domofon_id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var26 string
		templ_7745c5c3_Var26, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(props.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `client/web/views/components/singleDomofonCard.templ`, Line: 121, Col: 93}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var26))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 33, "\"> <button type=\"submit\" class=\"control-button delete-button\" title=\"Удалить домофон\"><svg width=\"16\" height=\"16\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\"><path d=\"M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2\"></path></svg></button></form></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func truncateText(text string, maxLength int) string {
	if len(text) > maxLength {
		return text[:maxLength] + "..."
	}
	return text
}

func doorIcon(status bool) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var27 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var27 == nil {
			templ_7745c5c3_Var27 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if status == false {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 34, "<img width=\"16px\" height=\"16px\" class=\"dooricon\" src=\"/client/web/static/icons/door-inactive.svg\" alt=\"dooractive ico\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 35, "<img width=\"16px\" height=\"16px\" class=\"dooricon\" src=\"/client/web/static/icons/door-active.svg\" alt=\"doorinactive ico\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return nil
	})
}

func phoneIcon(active bool) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var28 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var28 == nil {
			templ_7745c5c3_Var28 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 36, "<svg class=\"phone-icon\" width=\"16\" height=\"16\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if active {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 37, "<path d=\"M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z\" fill=\"#4CAF50\" stroke=\"#4CAF50\"></path>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 38, "<path d=\"M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z\" fill=\"none\" stroke=\"#9e9e9e\"></path>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 39, "</svg>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func ControlButtonStyle() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var29 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var29 == nil {
			templ_7745c5c3_Var29 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 40, "<style>\r\n        .power-controls {\r\n        display: flex;\r\n        gap: 6px;\r\n        align-items: center;\r\n    }\r\n        .delete-form {\r\n        margin-left: auto;\r\n    }\r\n        .delete-button {\r\n        background-color: #ff5252; !important;\r\n        color: white;\r\n        padding: 8px;\r\n        width: auto;\r\n    }\r\n        .delete-button:hover {\r\n        background-color: #d13c3c; !important;\r\n    }\r\n        .domofon-controls {\r\n            display: flex;\r\n            flex-direction: column;\r\n            gap: 12px;\r\n            width: 160px;\r\n        }\r\n    \r\n        .control-form {\r\n            display: flex;\r\n            flex-direction: column;\r\n            gap: 8px;\r\n        }\r\n    \r\n        .form-group {\r\n            display: flex;\r\n            flex-direction: column;\r\n            gap: 4px;\r\n        }\r\n    \r\n        .form-group label {\r\n            font-size: 0.75rem;\r\n            color: #555;\r\n        }\r\n    \r\n        .apartment-input {\r\n            padding: 6px 8px;\r\n            border: 1px solid #ddd;\r\n            border-radius: 4px;\r\n            font-size: 0.85rem;\r\n        }\r\n    \r\n        .control-button {\r\n            display: flex;\r\n            align-items: center;\r\n            justify-content: center;\r\n            gap: 6px;\r\n            padding: 8px 12px;\r\n            border: none;\r\n            border-radius: 4px;\r\n            font-size: 0.85rem;\r\n            cursor: pointer;\r\n            transition: background-color 0.2s;\r\n            \r\n        }\r\n    \r\n        .call-button {\r\n            background-color: #4CAF50;\r\n            color: white;\r\n        }\r\n    \r\n        .call-button:hover {\r\n            background-color: #3e8e41;\r\n        }\r\n    \r\n        .open-button {\r\n            background-color: #2196F3;\r\n            color: white;\r\n        }\r\n    \r\n        .open-button:hover {\r\n            background-color: #0b7dda;\r\n        }\r\n    \r\n        .power-buttonOFF {\r\n            background-color: #4caf50;\r\n            color: white;\r\n        }\r\n\r\n        .power-buttonON {\r\n            background-color: #ff5252;\r\n            color: white;\r\n        }\r\n    \r\n        .power-buttonOFF:hover {\r\n            background-color: #3e8e41;\r\n        }\r\n        .power-buttonON:hover {\r\n            background-color: #d13c3c;\r\n        }\r\n    </style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func SingleDomofonCardStyle() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var30 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var30 == nil {
			templ_7745c5c3_Var30 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 41, "<style>\r\n    .domofon-card-container {\r\n        display: flex;\r\n        height: 350px;\r\n        flex-direction: row;\r\n        align-items: flex-start;\r\n        gap: 20px;\r\n        margin-bottom: 24px;\r\n        background: #ffffff;\r\n        border-radius: 14px;\r\n        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);\r\n        padding: 20px;\r\n        transition: all 0.3s ease;\r\n        border: 1px solid #e0e0e0;\r\n    }\r\n\r\n    .door-status--wrapper {\r\n        display: flex;\r\n        align-items: center;\r\n        flex-direction: row;\r\n        gap: 8px;\r\n    }\r\n\r\n    .dooricon {\r\n        display: flex;\r\n        align-items: center;\r\n        width: 20px;\r\n        height: 20px;\r\n    }\r\n\r\n    .bell-and-door--wrapper {\r\n        display: flex;\r\n        align-items: center;\r\n        flex-direction: row;\r\n        text-align: center;\r\n        gap: 12px;\r\n    }\r\n\r\n    .phone-icon {\r\n        vertical-align: middle;\r\n        margin-right: 8px;\r\n        width: 20px;\r\n        height: 20px;\r\n    }\r\n\r\n    .call-indicator.calling .phone-icon {\r\n        animation: pulse 1.5s infinite;\r\n    }\r\n\r\n    .domofon-card {\r\n        background: white;\r\n        border-radius: 12px;\r\n        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);\r\n        padding: 20px;\r\n        width: 350px;\r\n        height: 100%;\r\n        display: flex;\r\n        flex-direction: column;\r\n        transition: all 0.2s ease;\r\n        border: 1px solid transparent;\r\n    }\r\n\r\n    .domofon-card.inactive {\r\n        border: 1px solid #ff5252;\r\n        position: relative;\r\n    }\r\n\r\n    .domofon-card.active {\r\n        border: 1px solid #4caf50;\r\n        position: relative;\r\n    }\r\n\r\n    .domofon-card.active::after {\r\n        content: \"\";\r\n        position: absolute;\r\n        top: 0;\r\n        left: 0;\r\n        right: 0;\r\n        bottom: 0;\r\n        border-radius: 11px;\r\n        border: 1px solid #4caf50;\r\n        pointer-events: none;\r\n    }\r\n\r\n    .domofon-card.inactive::after {\r\n        content: \"\";\r\n        position: absolute;\r\n        top: 0;\r\n        left: 0;\r\n        right: 0;\r\n        bottom: 0;\r\n        border-radius: 11px;\r\n        border: 1px solid #ff5252;\r\n        pointer-events: none;\r\n    }\r\n\r\n    .card-header {\r\n        display: flex;\r\n        justify-content: space-between;\r\n        align-items: center;\r\n        margin-bottom: 16px;\r\n        padding-bottom: 12px;\r\n        border-bottom: 1px solid #f0f0f0;\r\n    }\r\n\r\n    .domofon-id {\r\n        font-weight: bold;\r\n        color: #555;\r\n        font-size: 1.2rem;\r\n    }\r\n\r\n    .status-badge {\r\n        font-size: 0.9rem;\r\n        padding: 4px 10px;\r\n        border-radius: 14px;\r\n        font-weight: 500;\r\n        text-transform: uppercase;\r\n    }\r\n\r\n    .status-badge.active {\r\n        background: #e3fcef;\r\n        color: #008a45;\r\n    }\r\n\r\n    .status-badge.inactive {\r\n        background: #ffebee;\r\n        color: #d32f2f;\r\n    }\r\n\r\n    .status-badge.offline {\r\n        background: #f5f5f5;\r\n        color: #757575;\r\n    }\r\n\r\n    .door-status {\r\n        font-size: 1rem;\r\n    }\r\n\r\n    .door-status.open {\r\n        color: #4caf50;\r\n    }\r\n\r\n    .door-status.closed {\r\n        color: #f44336;\r\n    }\r\n\r\n    .card-body {\r\n        margin-bottom: 16px;\r\n        flex-grow: 1;\r\n        overflow: hidden;\r\n    }\r\n\r\n    .mac-address,\r\n    .address {\r\n        display: flex;\r\n        align-items: center;\r\n        gap: 10px;\r\n        font-size: 1rem;\r\n        margin-bottom: 12px;\r\n    }\r\n\r\n    .address {\r\n        margin-top: 8px;\r\n    }\r\n\r\n    .label {\r\n        color: #757575;\r\n        font-weight: 500;\r\n        font-size: 0.95rem;\r\n    }\r\n\r\n    .value {\r\n        white-space: nowrap;\r\n        overflow: hidden;\r\n        text-overflow: ellipsis;\r\n        font-size: 1rem;\r\n    }\r\n\r\n    .card-footer {\r\n        display: flex;\r\n        justify-content: space-between;\r\n        align-items: center;\r\n        font-size: 0.95rem;\r\n        color: #757575;\r\n    }\r\n\r\n    .apartments {\r\n        display: flex;\r\n        align-items: center;\r\n        gap: 8px;\r\n    }\r\n\r\n    .icon {\r\n        margin-right: 6px;\r\n    }\r\n</style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
