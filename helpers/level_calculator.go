package helpers

var levelValues = []uint64{
	30000,
	100000,
	210000,
	360000,
	550000,
	780000,
	1050000,
	1360000,
	1710000,
	2100000,
	2530000,
	3000000,
	3510000,
	4060000,
	4650000,
	5280000,
	5950000,
	6660000,
	7410000,
	8200000,
	9030000,
	9900000,
	10810000,
	11760000,
	12750000,
	13780000,
	14850000,
	15960000,
	17110000,
	18300000,
	19530000,
	20800000,
	22110000,
	23460000,
	24850000,
	26280000,
	27750000,
	29260000,
	30810000,
	32400000,
	34030000,
	35700000,
	37410000,
	39160000,
	40950000,
	42780000,
	44650000,
	46560000,
	48510000,
	50500000,
	52530000,
	54600000,
	56710000,
	58860000,
	61050000,
	63280000,
	65550000,
	67860000,
	70210001,
	72600001,
	75030002,
	77500003,
	80010006,
	82560010,
	85150019,
	87780034,
	90450061,
	93160110,
	95910198,
	98700357,
	101530643,
	104401157,
	107312082,
	110263748,
	113256747,
	116292144,
	119371859,
	122499346,
	125680824,
	128927482,
	132259468,
	135713043,
	139353477,
	143298259,
	147758866,
	153115959,
	160054726,
	169808506,
	184597311,
	208417160,
	248460887,
	317675597,
	439366075,
	655480935,
	1041527682,
	1733419828,
	2975801691,
	5209033044,
	9225761479,
	99999999999,
	99999999999,
	999999999999,
	999999999999,
	999999999999,
	999999999999,
	999999999999,
	999999999999,
	999999999999,
}

func GetLevelFromScore(score uint64) uint64 {
	var a, i uint64 = 0, 0

	for a+levelValues[i] < score {
		i++
		a += levelValues[i]
	}

	return i + 1 + (score-a)/levelValues[i+1]
}