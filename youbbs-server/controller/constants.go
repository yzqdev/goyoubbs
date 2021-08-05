package controller

import "goyoubbs/util"

func GetAppHome(prefix string) string {
	dir, _ := util.Dir()
	return dir + "/youbbs/static" + prefix

}
