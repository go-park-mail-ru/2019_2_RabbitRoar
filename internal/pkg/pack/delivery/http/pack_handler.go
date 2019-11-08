package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type handler struct {}

func NewPackHandler(
	e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
) {
	handler := handler {}

	group := e.Group("/pack", authMiddleware)

	group.GET("/:id", handler.byID)
}

func (h *handler) byID(ctx echo.Context) error {
	var (
		id, err = strconv.Atoi(ctx.Param("id"))
	)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "invalid pack id",
			Internal: err,
		}
	}
	if id == 0 {
		return ctx.JSONBlob(http.StatusOK, []byte(packMock0))
	} else if id == 1 {
		return ctx.JSONBlob(http.StatusOK, []byte(packMock0))
	}
	return echo.ErrNotFound
}

var packMock0 = `
{
   "id":0,
   "name":"Пак для лоуреатов",
   "img":"",
   "rating":50,
   "author":"AnitaKanita",
   "authorId":0,
   "questions":{
      "..ты..":[
         {
            "id":0,
            "text":"символом M в римской записи записывают это число",
            "media":"string",
            "answer":"тысяча",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"цвет, соответствующий оптическому диапазону длин волн 570—590 нм",
            "media":"string",
            "answer":"жёлтый",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"сериал про прапара Шматко",
            "media":"string",
            "answer":"солдаты",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"старинное название башенных или больших комнатных часов с набором настроенных колоколов, издающих бой в определённой мелодической последовательности",
            "media":"string",
            "answer":"куранты",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"денежная валюта Польши",
            "media":"string",
            "answer":"злотый",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "день Х":[
         {
            "id":0,
            "text":"Перевод на английский фразы «Да пребудет с тобой Сила» подскажет, когда празднуют день «Звёздных войн»",
            "media":"string",
            "answer":"4 мая",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"День программиста — праздник, отмечаемый в 256-й день года. Чаще всего праздник выпадает на 13 сентября, но иногда - именно на этот день.",
            "media":"string",
            "answer":"12 сентября",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"День холостяков — китайский современный праздник, посвящённый людям, не состоящим в браке. Получил своё название из-за того, что дата проведения символизирует не состоящих в паре людей. А ещё в этот день вышел Skyrim.",
            "media":"string",
            "answer":"11 ноября",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"День квадратного корня — неофициальный праздник, отмечаемый в день, когда и число, и порядковый номер месяца являются квадратными корнями из двух последних цифр года. Например, 2 февраля 2004 года. Когда мы отпразднуем в следующий раз?",
            "media":"string",
            "answer":"5 мая 2025",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"День числа пи (3,1415926...) — неофициальный праздник, который отмечается любителями математики 14 марта в 1:59:26. Некоторые также празднуют 22 июля. Почему?",
            "media":"string",
            "answer":"22/7 = 3.14",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "история успеха":[
         {
            "id":0,
            "text":"Томас Вудро Вильсон, 28-й президент США сказал: 'На подготовку 10-минутной речи мне нужна неделя; на 15-минутную - три дня; на получасовую - два дня; а часовую речь я могу [...]'",
            "media":"string",
            "answer":"произнести хоть сейчас",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"Каким мифическим существом называют частную стартап-компанию стоимостью более 1 миллиарда долларов? Термин был придуман в 2013 году для представления статистической редкости таких успешных предприятий.",
            "media":"string",
            "answer":"единорог",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"Самая продаваемая игровая консоль",
            "media":"string",
            "answer":"PlayStation 2",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"в прошлом крупнейшая частная российская авиакомпания, прекратившая свою операционную деятельность в октябре 2015 года",
            "media":"string",
            "answer":"трансаэро",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"Кодовые имена модельного ряда этой крупной американской компании можно сложить в 'S3XY'. Назовите компанию.",
            "media":"string",
            "answer":"Tesla",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "статья Lurkmore":[
         {
            "id":0,
            "text":"эффективная система техник и умений персонифицированного нагибания путём доставки пиздюлей кулаками в тушку и верхнюю голову оппонента",
            "media":"string",
            "answer":"бокс",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"злыдень, нацист, ефрейтор-полководец, носитель усов «Бродяга», поц и шлимазл, икона для фошыстов, потомок евреев и негров, талантливый оратор, выдающийся художник первой половины XX в., мамзер, пейсатель, звезда ютуба и просто тонкая натура",
            "media":"string",
            "answer":"Гитлер",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"межрасовая тоталитарная секта экологических экстремистов, состоящая из нескольких десятков обособленных контор, координирующих деятельность на ежегодных сходках директоров",
            "media":"string",
            "answer":"Гринпис",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"Про это есть порно. Никаких исключений",
            "media":"string",
            "answer":"Правило 34",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"Серия суровых 2d-комбо-файтингов с былинным сюжетом за авторством южноафриканского расово японского художника карандашами по бумаге и митолизда пальцами на гитаре Дайсуке Исиватари сотоварищи. Отличается от прочих представителей жанра повышенными концентрациями анимешности и, что самое ГЛАВНОЕ, рок-н-ролла и метала с прилагающимися отсылками к ним.",
            "media":"string",
            "answer":"Guilty Gear",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "международные меры":[
         {
            "id":0,
            "text":"так называют коллективные или односторонние принудительные меры, применяемые государствами или международными организациями к государству, которое нарушило нормы международного права",
            "media":"string",
            "answer":"санкции",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"английская версия названия 'камень-ножницы-бумага' имеет именно такой порядок компонентов",
            "media":"string",
            "answer":"камень-бумага-ножницы",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"лишь 19 из 28 стран Евросоюза пользуются им официально и полноценно",
            "media":"string",
            "answer":"евро",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"в грамме помещается столько карат",
            "media":"string",
            "answer":"5",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"если вы закажете пинту пива в Англии, то вам принесут стакан примерно такого объёма в литрах",
            "media":"string",
            "answer":"0.5",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ]
   },
   "tags":"tags string"
}
`
var packMock1 = `
{
   "id":1,
   "name":"Об играх",
   "img":"",
   "rating":75,
   "author":"EgosKekos",
   "authorId":1,
   "questions":{
      "Шутерки":[
         {
            "id":0,
            "text":"Модом к какой игре изначально была CS?",
            "media":"string",
            "answer":"Half-Life",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"Кто является компанией издателем небезызвестной Call of Duty",
            "media":"string",
            "answer":"Activision",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"Первая игра из серии Battlefield",
            "media":"string",
            "answer":"Battlefield 1942",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"Эта игра соединившая в себе черты рпг и шутера, известна своим генератором оружия",
            "media":"string",
            "answer":"Borderlands",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"Стелс шутер от 3го лица, симулятор отстрела яиц фрицев",
            "media":"string",
            "answer":"Sniper Elite",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "РПГ":[
         {
            "id":0,
            "text":"RPG разработанная канадской студией BioWare в 2009 г, как они сами ее называли «темное героическое фэнтези»",
            "media":"string",
            "answer":"Dragon Age: Origins",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"Назови все ведьмачьи знаки из серии игр Ведьмак (их пять)",
            "media":"string",
            "answer":"Аард Ирден Игни Квен Аксий",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"Слабо назвать все номерные части The Elder Scrolls?",
            "media":"string",
            "answer":"Arena Daggerfall Morrowind Oblivion Skyrim",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"BioWare анонсировала эту RPG на выставке Electronic Entertainment Expo 2015(E3)",
            "media":"string",
            "answer":"Mass Effect Andromeda",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"Первая игра позволявшая вступить в однополый брак. Считается олдфагами лучшей в серии.",
            "media":"string",
            "answer":"Fallout 2",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "Гоночки":[
         {
            "id":0,
            "text":"Вторая часть безумной гоночной аркады где можно 'выстреливать' водителем из лобового стекла, стараясь забросить его как можно дальше. ",
            "media":"string",
            "answer":"FlatOut 2",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"Гоночная игра с открытым миром, где картой выступает вся Америка от Юбейсофт",
            "media":"string",
            "answer":"The Crew",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"Последняя часть игрой серии Colin McRae Rally, с надписью Colin McRae в названии",
            "media":"string",
            "answer":"Colin McRae: DiRT",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"Оттуда все узнали бело-синюю бэху",
            "media":"string",
            "answer":"Need for speed most wanted",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"Гоночный краш тест симулятор",
            "media":"string",
            "answer":"BeamNG",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "Сурвайвал":[
         {
            "id":0,
            "text":"Один убийца, 4 выжившых и куча генераторов",
            "media":"string",
            "answer":"Dead by Daylight",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"В этой игре вы работник психиатрической лечебницы, один из пациентов которой очень хочет видеть вас в качестве своей невесты ",
            "media":"string",
            "answer":"Outlast Whistleblower",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"В этой игре можно угрожать пустым пистолетом",
            "media":"string",
            "answer":"I Am Alive",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"Основной геймплей этой игры - гребля и сбор кокосов",
            "media":"string",
            "answer":"STRANDED DEEP",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"ИГРА очень похожая на королевскую битву или голодные игры, вышедшая 8 марта 2016 ",
            "media":"string",
            "answer":"The Culling",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ],
      "Отечественный игрострой":[
         {
            "id":0,
            "text":"Компьютерная игра, изобретённая в СССР Алексеем Пажитновым в 1984г",
            "media":"string",
            "answer":"Тетрис",
            "rating":10,
            "author":0,
            "tags":[
               "string"
            ]
         },
         {
            "id":1,
            "text":"Лучший авиасимулятор, из всех, когда-либо выпускавшихся на PC",
            "media":"string",
            "answer":"Ил-2 Штурмовик",
            "rating":11,
            "author":1,
            "tags":[
               "string"
            ]
         },
         {
            "id":2,
            "text":"Йо-хо-хо, и бутылка рому! Главный симулятор пиратских нескучных будней",
            "media":"string",
            "answer":"Корсары",
            "rating":12,
            "author":2,
            "tags":[
               "string"
            ]
         },
         {
            "id":3,
            "text":"Чики-брики и в дамки",
            "media":"string",
            "answer":"stalker",
            "rating":13,
            "author":3,
            "tags":[
               "string"
            ]
         },
         {
            "id":4,
            "text":"Часть одной из самых известных серий пошаговых стратегий, разработаная российской компанией Nival Interactive",
            "media":"string",
            "answer":"Heroes of Might and Magic 5",
            "rating":14,
            "author":4,
            "tags":[
               "string"
            ]
         }
      ]
   },
   "tags":"tags string"
}
`
