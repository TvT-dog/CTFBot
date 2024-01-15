package main

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"regexp"
	"strconv"
)

func ATMessageEventHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		ctx := context.Background()
		//匹配是否为获取所有比赛
		CtfGames := regexp.MustCompile(`/Games `)
		if CtfGames.MatchString(data.Content) {
			CtfGames_num := CtfGames.FindStringSubmatch(data.Content)
			if len(CtfGames_num) > 0 {
				ToMessage := &dto.MessageToCreate{
					Content: Matchs(),
				}
				processor.method.PostMessage(ctx, data.ChannelID, ToMessage)
			}
		}
		//获取单个比赛信息
		CtfGame := regexp.MustCompile(`/CTF (\d+)`)
		if CtfGame.MatchString(data.Content) {
			CtfGame_num := CtfGame.FindStringSubmatch(data.Content)
			if len(CtfGame_num) > 0 {
				num, _ := strconv.Atoi(CtfGame_num[1])
				ToMessage := &dto.MessageToCreate{
					Content: Match(num),
				}
				processor.method.PostMessage(ctx, data.ChannelID, ToMessage)
			}
		}
		//进行用户注册
		Register := regexp.MustCompile(`/Reg 比赛ID:(\d+) 比赛队伍名称：(\w+)`)
		if Register.MatchString(data.Content) {
			Register_num := Register.FindStringSubmatch(data.Content)
			if len(Register_num) >= 3 {
				matchID := Register_num[1]
				matchTeam := Register_num[2]
				ToMessage := &dto.MessageToCreate{
					MessageReference: &dto.MessageReference{
						MessageID: data.ID,
					},
					Content: ReTeamCompetition(matchTeam, data.Author.Username, data.Author.ID, matchID),
				}
				processor.method.PostMessage(ctx, data.ChannelID, ToMessage)
			}
		}
		//获取比赛报名信息
		GetGame := regexp.MustCompile(`/Info 比赛ID:(\d+)`)
		if GetGame.MatchString(data.Content) {
			Register_num := GetGame.FindStringSubmatch(data.Content)
			if len(Register_num) > 0 {
				matchID := Register_num[1]
				ToMessage := &dto.MessageToCreate{
					MessageReference: &dto.MessageReference{
						MessageID: data.ID,
					},
					Content: getComptitionTeam(matchID),
				}
				processor.method.PostMessage(ctx, data.ChannelID, ToMessage)
			}
		}

		return nil
	}
}
