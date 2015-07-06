//  Copyright (c) 2015 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package db

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/couchbase/go-couchbase"
	"github.com/couchbase/sync_gateway/base"
)

var suffixMapOSX = []uint16{95, 105, 129, 131, 240, 258, 274, 496, 506, 532, 643, 677, 952, 966, 1056, 1062,
	1283, 1313, 1327, 1455, 1461, 1479, 1680, 1698, 1708, 1710, 1724, 1801, 1819, 1835, 1989, 1991,
	2096, 2106, 2132, 2243, 2277, 2495, 2505, 2529, 2531, 2640, 2658, 2674, 2949, 2951, 2965, 3055,
	3061, 3079, 3280, 3298, 3308, 3310, 3324, 3456, 3462, 3683, 3713, 3727, 3802, 3836, 3992, 4147,
	4173, 4202, 4236, 4392, 4544, 4568, 4570, 4601, 4619, 4635, 4789, 4791, 4880, 4898, 4908, 4910,
	4924, 5014, 5020, 5038, 5184, 5349, 5351, 5365, 5417, 5423, 5587, 5752, 5766, 5843, 5877, 6144,
	6168, 6170, 6201, 6219, 6235, 6389, 6391, 6547, 6573, 6602, 6636, 6792, 6883, 6913, 6927, 7017,
	7023, 7187, 7352, 7366, 7414, 7420, 7438, 7584, 7749, 7751, 7765, 7840, 7858, 7874, 8009, 8011,
	8025, 8181, 8199, 8354, 8360, 8378, 8412, 8426, 8582, 8757, 8763, 8846, 8872, 9142, 9176, 9207,
	9233, 9397, 9541, 9559, 9575, 9604, 9628, 9630, 9794, 9885, 9915, 9921, 9939, 10047, 10073, 10292,
	10302, 10336, 10444, 10468, 10470, 10689, 10691, 10701, 10719, 10735, 10808, 10810, 10824, 10980, 10998, 11084,
	11114, 11120, 11138, 11249, 11251, 11265, 11487, 11517, 11523, 11652, 11666, 11943, 11977, 12044, 12068, 12070,
	12289, 12291, 12301, 12319, 12335, 12447, 12473, 12692, 12702, 12736, 12813, 12827, 12983, 13087, 13117, 13123,
	13252, 13266, 13484, 13514, 13520, 13538, 13649, 13651, 13665, 13940, 13958, 13974, 14005, 14029, 14031, 14195,
	14340, 14358, 14374, 14406, 14432, 14596, 14743, 14777, 14852, 14866, 15156, 15162, 15213, 15227, 15383, 15555,
	15561, 15579, 15608, 15610, 15624, 15780, 15798, 15889, 15891, 15901, 15919, 15935, 16006, 16032, 16196, 16343,
	16377, 16405, 16429, 16431, 16595, 16740, 16758, 16774, 16849, 16851, 16865, 17155, 17161, 17179, 17208, 17210,
	17224, 17380, 17398, 17556, 17562, 17613, 17627, 17783, 17892, 17902, 17936, 18153, 18167, 18216, 18222, 18386,
	18548, 18550, 18564, 18615, 18621, 18639, 18785, 18894, 18904, 18928, 18930, 19000, 19018, 19034, 19188, 19190,
	19345, 19369, 19371, 19403, 19437, 19593, 19746, 19772, 19857, 19863, 20095, 20105, 20129, 20131, 20240, 20258,
	20274, 20496, 20506, 20532, 20643, 20677, 20952, 20966, 21056, 21062, 21283, 21313, 21327, 21455, 21461, 21479,
	21680, 21698, 21708, 21710, 21724, 21801, 21819, 21835, 21989, 21991, 22096, 22106, 22132, 22243, 22277, 22495,
	22505, 22529, 22531, 22640, 22658, 22674, 22949, 22951, 22965, 23055, 23061, 23079, 23280, 23298, 23308, 23310,
	23324, 23456, 23462, 23683, 23713, 23727, 23802, 23836, 23992, 24147, 24173, 24202, 24236, 24392, 24544, 24568,
	24570, 24601, 24619, 24635, 24789, 24791, 24880, 24898, 24908, 24910, 24924, 25014, 25020, 25038, 25184, 25349,
	25351, 25365, 25417, 25423, 25587, 25752, 25766, 25843, 25877, 26144, 26168, 26170, 26201, 26219, 26235, 26389,
	26391, 26547, 26573, 26602, 26636, 26792, 26883, 26913, 26927, 27017, 27023, 27187, 27352, 27366, 27414, 27420,
	27438, 27584, 27749, 27751, 27765, 27840, 27858, 27874, 28009, 28011, 28025, 28181, 28199, 28354, 28360, 28378,
	28412, 28426, 28582, 28757, 28763, 28846, 28872, 29142, 29176, 29207, 29233, 29397, 29541, 29559, 29575, 29604,
	29628, 29630, 29794, 29885, 29915, 29921, 29939, 30047, 30073, 30292, 30302, 30336, 30444, 30468, 30470, 30689,
	30691, 30701, 30719, 30735, 30808, 30810, 30824, 30980, 30998, 31084, 31114, 31120, 31138, 31249, 31251, 31265,
	31487, 31517, 31523, 31652, 31666, 31943, 31977, 32044, 32068, 32070, 32289, 32291, 32301, 32319, 32335, 32447,
	32473, 32692, 32702, 32736, 32813, 32827, 32983, 33087, 33117, 33123, 33252, 33266, 33484, 33514, 33520, 33538,
	33649, 33651, 33665, 33940, 33958, 33974, 34005, 34029, 34031, 34195, 34340, 34358, 34374, 34406, 34432, 34596,
	34743, 34777, 34852, 34866, 35156, 35162, 35213, 35227, 35383, 35555, 35561, 35579, 35608, 35610, 35624, 35780,
	35798, 35889, 35891, 35901, 35919, 35935, 36006, 36032, 36196, 36343, 36377, 36405, 36429, 36431, 36595, 36740,
	36758, 36774, 36849, 36851, 36865, 37155, 37161, 37179, 37208, 37210, 37224, 37380, 37398, 37556, 37562, 37613,
	37627, 37783, 37892, 37902, 37936, 38153, 38167, 38216, 38222, 38386, 38548, 38550, 38564, 38615, 38621, 38639,
	38785, 38894, 38904, 38928, 38930, 39000, 39018, 39034, 39188, 39190, 39345, 39369, 39371, 39403, 39437, 39593,
	39746, 39772, 39857, 39863, 40095, 40105, 40129, 40131, 40240, 40258, 40274, 40496, 40506, 40532, 40643, 40677,
	40952, 40966, 41056, 41062, 41283, 41313, 41327, 41455, 41461, 41479, 41680, 41698, 41708, 41710, 41724, 41801,
	41819, 41835, 41989, 41991, 42096, 42106, 42132, 42243, 42277, 42495, 42505, 42529, 42531, 42640, 42658, 42674,
	42949, 42951, 42965, 43055, 43061, 43079, 43280, 43298, 43308, 43310, 43324, 43456, 43462, 43683, 43713, 43727,
	43802, 43836, 43992, 44147, 44173, 44202, 44236, 44392, 44544, 44568, 44570, 44601, 44619, 44635, 44789, 44791,
	44880, 44898, 44908, 44910, 44924, 45014, 45020, 45038, 45184, 45349, 45351, 45365, 45417, 45423, 45587, 45752,
	45766, 45843, 45877, 46144, 46168, 46170, 46201, 46219, 46235, 46389, 46391, 46547, 46573, 46602, 46636, 46792,
	46883, 46913, 46927, 47017, 47023, 47187, 47352, 47366, 47414, 47420, 47438, 47584, 47749, 47751, 47765, 47840,
	47858, 47874, 48009, 48011, 48025, 48181, 48199, 48354, 48360, 48378, 48412, 48426, 48582, 48757, 48763, 48846,
	48872, 49142, 49176, 49207, 49233, 49397, 49541, 49559, 49575, 49604, 49628, 49630, 49794, 49885, 49915, 49921,
	49939, 50047, 50073, 50292, 50302, 50336, 50444, 50468, 50470, 50689, 50691, 50701, 50719, 50735, 50808, 50810,
	50824, 50980, 50998, 51084, 51114, 51120, 51138, 51249, 51251, 51265, 51487, 51517, 51523, 51652, 51666, 51943,
	51977, 52044, 52068, 52070, 52289, 52291, 52301, 52319, 52335, 52447, 52473, 52692, 52702, 52736, 52813, 52827,
	52983, 53087, 53117, 53123, 53252, 53266, 53484, 53514, 53520, 53538, 53649, 53651, 53665, 53940, 53958, 53974,
	54005, 54029, 54031, 54195, 54340, 54358, 54374, 54406, 54432, 54596, 54743, 54777, 54852, 54866, 55156, 55162,
	55213, 55227, 55383, 55555, 55561, 55579, 55608, 55610, 55624, 55780, 55798, 55889, 55891, 55901, 55919, 55935,
	56006, 56032, 56196, 56343, 56377, 56405, 56429, 56431, 56595, 56740, 56758, 56774, 56849, 56851, 56865, 57155,
	57161, 57179, 57208, 57210, 57224, 57380, 57398, 57556, 57562, 57613, 57627, 57783, 57892, 57902, 57936, 58153,
	58167, 58216, 58222, 58386, 58548, 58550, 58564, 58615, 58621, 58639, 58785, 58894, 58904, 58928, 58930, 59000,
	59018, 59034, 59188, 59190, 59345, 59369, 59371, 59403, 59437, 59593, 59746, 59772, 59857, 59863, 60095, 60105,
	60129, 60131, 60240, 60258, 60274, 60496, 60506, 60532, 60643, 60677, 60952, 60966, 61056, 61062, 61283, 61313,
	61327, 61455, 61461, 61479, 61680, 61698, 61708, 61710, 61724, 61801, 61819, 61835, 61989, 61991, 62096, 62106,
	62132, 62243, 62277, 62495, 62505, 62529, 62531, 62640, 62658, 62674, 62949, 62951, 62965, 63055, 63061, 63079,
	63280, 63298, 63308, 63310, 63324, 63456, 63462, 63683, 63713, 63727, 63802, 63836, 63992, 64147, 64173, 64202,
	64236, 64392, 64544, 64568, 64570, 64601, 64619, 64635, 64789, 64791, 64880, 64898, 64908, 64910, 64924, 65014,
}

var suffixMapRemote = []uint32{609, 799, 3428, 5688, 5718, 6539, 10876, 11172, 11600, 11790, 12353, 12421, 14063, 14681, 14711, 15967, 17242, 17530, 19287, 19317, 19465, 20934, 21030, 21742, 22211, 22381, 22563, 24121, 24653, 25825, 27290, 27300, 27472, 29255, 29527, 30039, 33218, 33388, 35128, 36299, 36309, 40226, 40554, 42893, 42903, 43007, 43197, 43775, 45337, 45445, 46086, 46116, 46664, 47812, 47982, 48043, 48731, 49947, 59738, 63848, 66959, 69098, 69108, 70364, 70416, 70586, 72841, 73145, 73637, 75275, 75497, 75507, 76054, 76726, 77950, 78091, 78101, 78673, 79805, 79995, 80529, 83698, 83708, 85438, 86619, 86789, 91252, 91520, 92073, 92691, 92701, 93977, 94343, 94431, 96866, 97162, 97610, 97780, 98933, 99037, 99745, 100899, 100909, 105818, 105988, 109268, 110004, 110194, 110776, 111890, 111900, 113225, 113557, 114811, 114981, 115085, 115115, 115667, 116334, 116446, 118261, 118483, 118513, 120146, 120634, 121842, 123367, 123415, 123585, 124953, 125057, 125725, 126276, 126494, 126504, 128323, 128451, 139458, 141350, 141422, 142171, 142603, 142793, 143875, 144241, 144533, 146964, 147060, 147682, 147712, 148831, 149135, 149647, 150359, 153178, 155248, 156069, 159838, 160569, 163748, 165478, 166659, 171212, 171382, 171560, 172033, 172741, 173937, 174293, 174303, 174471, 176826, 177122, 177650, 178973, 179077, 179695, 179705, 183808, 183998, 186889, 186919, 189148, 190324, 190456, 192801, 192991, 193095, 193105, 193677, 195235, 195547, 196014, 196184, 196766, 197880, 197910, 198141, 198633, 199845, 200124, 200656, 201820, 203295, 203305, 203477, 204931, 205035, 205747, 206214, 206384, 206566, 208341, 208433, 210829, 215938, 219348, 229578, 230066, 230684, 230714, 231962, 233247, 233535, 234873, 235177, 235605, 235795, 236356, 236424, 238203, 238393, 238571, 240279, 243058, 245368, 246149, 249888, 249918, 251270, 251492, 251502, 252051, 252723, 253955, 254361, 254413, 254583, 256844, 257140, 257632, 258881, 258911, 259015, 259185, 259767, 261332, 261440, 262083, 262113, 262661, 263817, 263987, 264223, 264551, 266896, 266906, 267002, 267192, 267770, 268853, 269157, 269625, 270449, 273668, 275558, 276779, 280204, 280394, 280576, 282921, 283025, 283757, 285285, 285315, 285467, 286134, 286646, 287830, 288061, 288683, 288713, 289965, 293928, 296839, 299068, 300956, 301052, 301720, 302273, 302491, 302501, 304143, 304631, 305847, 307362, 307410, 307580, 309237, 309545, 310729, 313498, 313508, 315638, 316419, 316589, 320089, 320119, 323338, 325008, 325198, 326229, 330814, 330984, 331080, 331110, 331662, 332331, 332443, 334001, 334191, 334773, 335895, 335905, 337220, 337552, 339375, 339407, 339597, 349618, 349788, 350296, 350306, 350474, 352823, 353127, 353655, 355217, 355387, 355565, 356036, 356744, 357932, 358163, 358611, 358781, 359867, 360244, 360536, 362961, 363065, 363687, 363717, 365355, 365427, 366174, 366606, 366796, 367870, 368021, 368753, 369925, 373968, 376879, 379028, 381372, 381400, 381590, 382153, 382621, 383857, 384263, 384481, 384511, 386946, 387042, 387730, 388813, 388983, 389087, 389117, 389665, 390409, 390599, 393628, 395488, 395518, 396739, 401639, 402418, 402588, 404728, 407499, 407509, 410142, 410630, 411846, 413363, 413411, 413581, 414957, 415053, 415721, 416272, 416490, 416500, 418327, 418455, 420000, 420190, 420772, 421894, 421904, 423221, 423553, 424815, 424985, 425081, 425111, 425663, 426330, 426442, 428265, 428487, 428517, 431009, 431199, 432228, 434088, 434118, 437339, 441216, 441386, 441564, 442037, 442745, 443933, 444297, 444307, 444475, 446822, 447126, 447654, 448977, 449073, 449691, 449701, 458698, 458708, 462878, 467969, 468138, 471354, 471426, 472175, 472607, 472797, 473871, 474245, 474537, 476960, 477064, 477686, 477716, 478835, 479131, 479643, 481489, 481519, 482738, 484408, 484598, 487629, 490262, 490480, 490510, 492947, 493043, 493731, 495373, 495401, 495591, 496152, 496620, 497856, 498007, 498197, 498775, 499893, 499903, 501939, 504828, 508258, 510930, 511034, 511746, 512215, 512385, 512567, 514125, 514657, 515821, 517294, 517304, 517476, 519251, 519523, 520872, 521176, 521604, 521794, 522357, 522425, 524067, 524685, 524715, 525963, 527246, 527534, 529283, 529313, 529461, 538468, 540360, 540412, 540582, 542845, 543141, 543633, 545271, 545493, 545503, 546050, 546722, 547954, 548095, 548105, 548677, 549801, 549991, 551369, 552148, 554278, 557059, 558808, 558998, 561559, 562778, 564448, 567669, 570222, 570550, 572897, 572907, 573003, 573193, 573771, 575333, 575441, 576082, 576112, 576660, 577816, 577986, 578047, 578735, 579943, 582838, 587929, 588178, 591284, 591314, 591466, 592135, 592647, 593831, 594205, 594395, 594577, 596920, 597024, 597756, 598875, 599171, 599603, 599793, 600810, 600980, 601084, 601114, 601666, 602335, 602447, 604005, 604195, 604777, 605891, 605901, 607224, 607556, 609371, 609403, 609593, 611819, 611989, 614898, 614908, 618378, 628548, 630952, 631056, 631724, 632277, 632495, 632505, 634147, 634635, 635843, 637366, 637414, 637584, 639233, 639541, 641249, 642068, 644358, 647179, 648928, 650240, 650532, 652965, 653061, 653683, 653713, 655351, 655423, 656170, 656602, 656792, 657874, 658025, 658757, 659921, 660292, 660302, 660470, 662827, 663123, 663651, 665213, 665383, 665561, 666032, 666740, 667936, 668167, 668615, 668785, 669863, 671479, 672658, 674568, 677749, 681234, 681546, 682015, 682185, 682767, 683881, 683911, 684325, 684457, 686800, 686990, 687094, 687104, 687676, 688955, 689051, 689723, 692888, 692918, 697809, 697999, 698058, 700062, 700680, 700710, 701966, 703243, 703531, 704877, 705173, 705601, 705791, 706352, 706420, 708207, 708397, 708575, 711689, 711719, 712538, 714608, 714798, 717429, 721129, 722298, 722308, 724038, 727219, 727389, 730120, 730652, 731824, 733291, 733301, 733473, 734935, 735031, 735743, 736210, 736380, 736562, 738345, 738437, 748628, 751336, 751444, 752087, 752117, 752665, 753813, 753983, 754227, 754555, 756892, 756902, 757006, 757196, 757774, 758857, 759153, 759621, 761274, 761496, 761506, 762055, 762727, 763951, 764365, 764417, 764587, 766840, 767144, 767636, 768885, 768915, 769011, 769181, 769763, 772958, 777849, 778018, 778188, 780342, 780430, 782867, 783163, 783611, 783781, 785253, 785521, 786072, 786690, 786700, 787976, 788127, 788655, 789823, 791439, 792618, 792788, 794528, 797699, 797709, 801448, 802669, 804559, 807778, 810333, 810441, 812816, 812986, 813082, 813112, 813660, 815222, 815550, 816003, 816193, 816771, 817897, 817907, 818156, 818624, 819852, 820271, 820493, 820503, 822954, 823050, 823722, 825360, 825412, 825582, 826141, 826633, 827845, 828014, 828184, 828766, 829880, 829910, 831278, 832059, 834369, 837148, 838889, 838919, 840963, 841067, 841685, 841715, 842246, 842534, 844176, 844604, 844794, 845872, 847357, 847425, 849202, 849392, 849570, 858579, 861828, 864939, 868349, 870821, 871125, 871657, 872294, 872304, 872476, 874034, 874746, 875930, 877215, 877385, 877567, 879340, 879432, 881768, 882549, 884679, 887458, 890013, 890183, 890761, 891887, 891917, 893232, 893540, 894806, 894996, 895092, 895102, 895670, 896323, 896451, 898276, 898494, 898504, 902969, 907878, 908029, 911245, 911537, 912064, 912686, 912716, 913960, 914354, 914426, 916871, 917175, 917607, 917797, 918924, 919020, 919752, 921297, 921307, 921475, 922126, 922654, 923822, 924216, 924386, 924564, 926933, 927037, 927745, 928866, 929162, 929610, 929780, 938619, 938789, 940081, 940111, 940663, 941815, 941985, 943330, 943442, 944894, 944904, 945000, 945190, 945772, 946221, 946553, 948374, 948406, 948596, 951088, 951118, 952339, 954009, 954199, 957228, 961728, 962499, 962509, 964639, 967418, 967588, 970053, 970721, 971957, 973272, 973490, 973500, 974846, 975142, 975630, 976363, 976411, 976581, 978236, 978544, 981868, 984979, 988299, 988309, 990861, 991165, 991617, 991787, 992344, 992436, 994074, 994696, 994706, 995970, 997255, 997527, 999290, 999300, 999472, 1001299, 1001309, 1002128, 1004218, 1004388, 1007039, 1008868, 1010290, 1010300, 1010472, 1012825, 1013121, 1013653, 1015211, 1015381, 1015563, 1016030, 1016742, 1017934, 1018165, 1018617, 1018787, 1019861, 1020242, 1020530, 1022967, 1023063, 1023681, 1023711, 1025353, 1025421, 1026172, 1026600}

var suffixMap = suffixMapOSX

func testIndexBucket() base.Bucket {
	/*
		bucket, err := ConnectToBucket(base.BucketSpec{
			Server:     "http://localhost:8091",
			BucketName: "channel_index"})
	*/

	/*
		bucket, err := ConnectToBucket(base.BucketSpec{
			Server:     "http://172.23.96.62:8091",
			BucketName: "channel_index"})
	*/

	bucket, err := ConnectToBucket(base.BucketSpec{
		Server:     "http://localhost:8091",
		BucketName: "channel_index"})

	if err != nil {
		log.Fatalf("Couldn't connect to bucket: %v", err)
	}
	return bucket
}

type channelIndexTest struct {
	numVbuckets   int         // Number of vbuckets
	indexBucket   base.Bucket // Index Bucket
	sequenceGap   int         // Max sequence gap within vbucket - random between 1 and sequenceGap
	lastSequences []uint64    // Last sequence per vbucket
	r             *rand.Rand  // seeded random number generator
	channelName   string      // Channel name
}

func NewChannelIndex(vbNum int, sequenceGap int, name string) *channelIndexTest {
	lastSeqs := make([]uint64, vbNum)

	index := &channelIndexTest{
		numVbuckets:   vbNum,
		indexBucket:   testIndexBucket(),
		sequenceGap:   sequenceGap,
		lastSequences: lastSeqs,
		r:             rand.New(rand.NewSource(42)),
		channelName:   name,
	}
	couchbase.PoolSize = 64
	return index
}

func NewChannelIndexForBucket(vbNum int, sequenceGap int, name string, bucket base.Bucket) *channelIndexTest {
	lastSeqs := make([]uint64, vbNum)

	index := &channelIndexTest{
		numVbuckets:   vbNum,
		indexBucket:   bucket,
		sequenceGap:   sequenceGap,
		lastSequences: lastSeqs,
		r:             rand.New(rand.NewSource(42)),
		channelName:   name,
	}
	couchbase.PoolSize = 64
	return index
}

func (c *channelIndexTest) seedData(format string) error {

	// Check if the data has already been loaded
	loadFlagDoc := fmt.Sprintf("seedComplete::%s", format)
	_, err := c.indexBucket.GetRaw(loadFlagDoc)
	if err == nil {
		return nil
	}

	vbucketBytes := make(map[int][]byte)

	// Populate index
	for i := 0; i < 100000; i++ {
		// Choose a vbucket at random
		vbNo := c.r.Intn(c.numVbuckets)
		_, ok := vbucketBytes[vbNo]
		if !ok {
			vbucketBytes[vbNo] = []byte("")
		}
		nextSequence := c.getNextSequenceBytes(vbNo)
		c.lastSequences[vbNo] = nextSequence
		vbucketBytes[vbNo] = append(vbucketBytes[vbNo], getSequenceAsBytes(nextSequence)...)
	}

	for vbNum, value := range vbucketBytes {
		_, err = c.indexBucket.AddRaw(c.getIndexDocName(vbNum), 0, value)
	}

	// Write flag doc
	_, err = c.indexBucket.AddRaw(loadFlagDoc, 0, []byte("complete"))
	if err != nil {
		log.Printf("Load error %v", err)
		return err
	}
	return nil
}

func (c *channelIndexTest) getNextSequenceBytes(vb int) uint64 {
	lastSequence := c.lastSequences[vb]
	gap := 1
	if c.sequenceGap > 0 {
		gap = 1 + c.r.Intn(c.sequenceGap)
	}
	nextSequence := lastSequence + uint64(gap)
	return nextSequence
}

func (c *channelIndexTest) addToCache(vb int) error {

	lastSequence := c.lastSequences[vb]
	gap := 1
	if c.sequenceGap > 0 {
		gap = 1 + c.r.Intn(c.sequenceGap)
	}
	nextSequence := lastSequence + uint64(gap)
	c.writeToCache(vb, nextSequence)
	c.lastSequences[vb] = nextSequence
	return nil
}

func (c *channelIndexTest) writeToCache(vb int, sequence uint64) error {
	docName := c.getIndexDocName(vb)
	sequenceBytes := getSequenceAsBytes(sequence)
	err := c.indexBucket.Append(docName, sequenceBytes)
	if err != nil {
		added, err := c.indexBucket.AddRaw(docName, 0, sequenceBytes)
		if err != nil || added == false {
			log.Printf("AddRaw also failed?! %s:%v", docName, err)
		}
	}
	return nil
}

func appendWrite(bucket base.Bucket, key string, vb int, sequence uint64) error {
	sequenceBytes := getSequenceAsBytes(sequence)
	err := bucket.Append(key, sequenceBytes)
	// TODO: assumes err means it hasn't been created yet
	if err != nil {
		added, err := bucket.AddRaw(key, 0, sequenceBytes)
		if err != nil || added == false {
			log.Printf("AddRaw also failed?! %s:%v", key, err)
		}
	}
	return nil
}

func getDocName(channelName string, vb int, blockNumber int) string {
	return fmt.Sprintf("_index::%s::%d::%05d", channelName, blockNumber, suffixMap[vb])
}

func (c *channelIndexTest) getIndexDocName(vb int) string {
	blockNumber := 0
	return getDocName(c.channelName, vb, blockNumber)
}

func (c *channelIndexTest) readIndexSingle() error {

	for i := 0; i < c.numVbuckets; i++ {
		key := c.getIndexDocName(i)
		body, err := c.indexBucket.GetRaw(key)
		if err != nil {
			log.Printf("Error retrieving for key %s: %s", key, err)
		} else {
			size := len(body)
			if size == 0 {
				return errors.New("Empty body")
			}
		}
	}
	return nil
}

func (c *channelIndexTest) readIndexSingleParallel() error {

	var wg sync.WaitGroup
	for i := 0; i < c.numVbuckets; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := c.getIndexDocName(i)
			body, err := c.indexBucket.GetRaw(key)
			if err != nil {
				log.Printf("Error retrieving for key %s: %s", key, err)
			} else {
				size := len(body)
				if size == 0 {
					log.Printf("zero size body found")
				}
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func (c *channelIndexTest) readIndexBulk() error {
	keys := make([]string, c.numVbuckets)
	for i := 0; i < c.numVbuckets; i++ {
		keys[i] = c.getIndexDocName(i)
	}
	couchbaseBucket, ok := c.indexBucket.(base.CouchbaseBucket)
	if !ok {
		log.Printf("Unable to convert to couchbase bucket")
		return errors.New("Unable to convert to couchbase bucket")
	}
	responses, err := couchbaseBucket.GetBulk(keys)
	if err != nil {
		return err
	}

	for _, response := range responses {
		body := response.Body
		// read last from body
		size := len(body)
		if size <= 0 {
			return errors.New(fmt.Sprintf("Empty body for response %v", response))
		}
	}
	return nil

}

func getSequenceAsBytes(sequence uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, sequence)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

// set up bucket connection

// define channel index (num vbuckets)

// write docs with random(?) vbucket distribution.
// include random gaps in sequence numbering within bucket

// benchmark
// init with vbucket size
//
// write n docs to channel index
// -
// get since 0
// write n docs to channel index
// get since previous

// bonus points:
// - check vbucket hash for docs

func BenchmarkChannelIndexSimpleGet(b *testing.B) {

	// num vbuckets
	vbCount := 1024
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	index.seedData("default")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := index.readIndexSingle()
		if err != nil {
			log.Printf("Error reading index single: %v", err)
		}
	}
}

func BenchmarkChannelIndexSimpleParallelGet(b *testing.B) {

	// num vbuckets
	vbCount := 1024
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	index.seedData("default")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := index.readIndexSingleParallel()
		if err != nil {
			log.Printf("Error reading index single: %v", err)
		}
	}
}

func BenchmarkChannelIndexBulkGet(b *testing.B) {

	// num vbuckets
	vbCount := 1024
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	index.seedData("default")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := index.readIndexBulk()
		if err != nil {
			log.Printf("Error reading index bulk: %v", err)
		}
	}
}

func BenchmarkChannelIndexPartitionReadSimple(b *testing.B) {

	// num vbuckets
	vbCount := 16
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	// Populate index with 100K sequences

	index.seedData("default")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := index.readIndexSingle()
		if err != nil {
			log.Printf("Error reading index single: %v", err)
		}
	}
}

func BenchmarkChannelIndexPartitionReadBulk(b *testing.B) {

	// num vbuckets
	vbCount := 16
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	index.seedData("default")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := index.readIndexBulk()
		if err != nil {
			log.Printf("Error reading index bulk: %v", err)
		}
	}
}

func BenchmarkMultiChannelIndexSimpleGet_10(b *testing.B) {
	MultiChannelIndexSimpleGet(b, 10)
}

func BenchmarkMultiChannelIndexSimpleGet_100(b *testing.B) {
	MultiChannelIndexSimpleGet(b, 100)
}

func BenchmarkMultiChannelIndexSimpleGet_1000(b *testing.B) {
	MultiChannelIndexSimpleGet(b, 1000)
}

func MultiChannelIndexSimpleGet(b *testing.B, numChannels int) {

	// num vbuckets
	vbCount := 1024

	bucket := testIndexBucket()
	indices := seedMultiChannelData(vbCount, bucket, numChannels)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for index := 0; index < numChannels; index++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				log.Printf("init %d", index)
				err := indices[index].readIndexSingle()
				log.Printf("done %d", index)
				if err != nil {
					log.Printf("Error reading index single: %v", err)
				}
			}(index)
		}
		wg.Wait()
	}
}

func seedMultiChannelData(vbCount int, bucket base.Bucket, numChannels int) []*channelIndexTest {

	indices := make([]*channelIndexTest, numChannels)
	// seed data
	for chanIndex := 0; chanIndex < numChannels; chanIndex++ {
		channelName := fmt.Sprintf("channel_%d", chanIndex)
		indices[chanIndex] = NewChannelIndexForBucket(vbCount, 0, channelName, bucket)
		indices[chanIndex].seedData(channelName)
	}

	log.Printf("Load complete")
	return indices
}
func BenchmarkMultiChannelIndexBulkGet_3(b *testing.B) {
	MultiChannelIndexBulkGet(b, 3)
}

func BenchmarkMultiChannelIndexBulkGet_10(b *testing.B) {
	MultiChannelIndexBulkGet(b, 10)
}

func BenchmarkMultiChannelIndexBulkGet_100(b *testing.B) {
	MultiChannelIndexBulkGet(b, 100)
}

func BenchmarkMultiChannelIndexBulkGet_1000(b *testing.B) {
	MultiChannelIndexBulkGet(b, 1000)
}

func MultiChannelIndexBulkGet(b *testing.B, numChannels int) {

	// num vbuckets
	vbCount := 1024

	bucket := testIndexBucket()
	indices := seedMultiChannelData(vbCount, bucket, numChannels)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for index := 0; index < numChannels; index++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				log.Printf("Calling bulk read for %d", index)
				err := indices[index].readIndexBulk()
				if err != nil {
					log.Printf("Error reading index single: %v", err)
				}
			}(index)
		}
		wg.Wait()
	}
}

func TestChannelIndexBulkGet10(t *testing.T) {

	// num vbuckets
	vbCount := 1024
	numChannels := 10
	bucket := testIndexBucket()
	indices := seedMultiChannelData(vbCount, bucket, numChannels)

	startTime := time.Now()

	var wg sync.WaitGroup
	for index := 0; index < numChannels; index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			err := indices[index].readIndexBulk()
			if err != nil {
				log.Printf("Error reading index single: %v", err)
			}
		}(index)
	}
	wg.Wait()

	log.Printf("test took %v", time.Since(startTime))
}

func TestChannelIndexSimpleReadSingle(t *testing.T) {

	log.Printf("Test single...")
	// num vbuckets
	vbCount := 1024
	index := NewChannelIndex(vbCount, 10, "basicChannel")

	// Populate index
	for i := 0; i < 5000; i++ {
		// Choose a vbucket
		vbNo := index.r.Intn(vbCount)
		index.addToCache(vbNo)
	}

	index.readIndexSingle()
}

func TestChannelIndexSimpleReadBulk(t *testing.T) {

	log.Printf("Test bulk...")
	// num vbuckets
	vbCount := 1024
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	index.readIndexBulk()
}

func TestChannelIndexPartitionReadSingle(t *testing.T) {

	log.Printf("Test single...")
	// num vbuckets
	vbCount := 16
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	// Populate index
	for i := 0; i < 5000; i++ {
		// Choose a vbucket
		vbNo := index.r.Intn(vbCount)
		index.addToCache(vbNo)
	}

	index.readIndexSingle()
}

func TestChannelIndexPartitionReadBulk(t *testing.T) {

	log.Printf("Test single...")
	// num vbuckets
	vbCount := 16
	index := NewChannelIndex(vbCount, 0, "basicChannel")

	// Populate index
	for i := 0; i < 5000; i++ {
		// Choose a vbucket
		vbNo := index.r.Intn(vbCount)
		index.addToCache(vbNo)
	}

	index.readIndexSingle()
}

func TestVbucket(t *testing.T) {

	index := NewChannelIndex(1024, 0, "basicChannel")
	counts := make(map[uint32]int)
	results := ""
	for i := 0; i < 10000000; i++ {
		docId := fmt.Sprintf("basicChannel::index::indexDoc_%07d", i)
		vbNo := index.indexBucket.VBHash(docId)
		if vbNo == 569 {
			results = fmt.Sprintf("%s, %d", results, i)
		}
		counts[vbNo] = counts[vbNo] + 1
		if counts[vbNo] > 1024 {
			log.Printf("your winner %d, at %d:", vbNo, i)
			break
		}
	}
	log.Printf(results)
	log.Printf("done")

}

func verifyVBMapping(bucket base.Bucket, channelName string) error {

	channelVbNo := uint32(0)

	for i := 0; i < 1024; i++ {
		docId := fmt.Sprintf("_index::%s::%d::%05d", channelName, 1, suffixMap[i])
		vbNo := bucket.VBHash(docId)
		if channelVbNo == 0 {
			channelVbNo = vbNo
			log.Printf("channel %s gets vb %d", channelName, channelVbNo)
		}
		if vbNo != channelVbNo {
			log.Println("Sad trombone - no match")
			return errors.New("vb numbers don't match")
		}
	}
	return nil
}

func TestChannelVbucketMappings(t *testing.T) {

	index := NewChannelIndex(1024, 0, "basicChannel")

	err := verifyVBMapping(index.indexBucket, "foo")
	assertTrue(t, err == nil, "inconsistent hash")

	err = verifyVBMapping(index.indexBucket, "SomeVeryLongChannelNameInCaseLengthIsSomehowAFactor")
	assertTrue(t, err == nil, "inconsistent hash")
	err = verifyVBMapping(index.indexBucket, "Punc-tu-@-tio-n")
	assertTrue(t, err == nil, "inconsistent hash")
	err = verifyVBMapping(index.indexBucket, "more::punc::tu::a::tion")
	assertTrue(t, err == nil, "inconsistent hash")

	err = verifyVBMapping(index.indexBucket, "1")
	assertTrue(t, err == nil, "inconsistent hash")

	log.Printf("checks out")
}