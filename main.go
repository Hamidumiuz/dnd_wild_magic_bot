package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	// Получаем токен из переменной окружения
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN не установлен")
	}

	// Создаем бота
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatalf("bot init failed: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Print("start bot")
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	// Обработка команды /start
	if update.Message.Text == "/start" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Икбол будет забанен)))",
		})
		return
	}

	// Бросок кубиков и поиск результата
	roll := rollDice(1, 100)
	result := wildMagicTable.Find(roll)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Бросок к100: %d\n%s", roll, result),
		ParseMode: models.ParseModeMarkdownV1,
	})
}

var rander = rand.New(rand.NewSource(time.Now().UnixNano()))

func rollDice(numDice, numSides int) int {
	total := 0
	for i := 0; i < numDice; i++ {
		total += rander.Intn(numSides) + 1
	}
	return total
}

type DiceRange struct {
	Min int
	Max int
}

type DiceTableRow struct {
	Range DiceRange
	Value func() string
}

type DiceTable []DiceTableRow

func (dt DiceTable) Len() int           { return len(dt) }
func (dt DiceTable) Less(i, j int) bool { return dt[i].Range.Min < dt[j].Range.Min }
func (dt DiceTable) Swap(i, j int)      { dt[i], dt[j] = dt[j], dt[i] }

func (dt DiceTable) Find(roll int) string {
	idx := sort.Search(len(dt), func(i int) bool {
		return dt[i].Range.Max >= roll
	})

	if idx < len(dt) && roll >= dt[idx].Range.Min && roll <= dt[idx].Range.Max {
		return dt[idx].Value()
	}
	return "Значение не найдено"
}

var wildMagicTable = DiceTable{
	{Range: DiceRange{Min: 1, Max: 4}, Value: func() string {
		return "Совершайте бросок по этой таблице в начале каждого своего хода в течение следующей минуты, игнорируя результат «01–04»."
	}},
	{Range: DiceRange{Min: 5, Max: 8}, Value: func() string {
		idx := rander.Intn(len(effects_5_8))
		return fmt.Sprintf("Дружелюбное существо появляется в пределах 60 футов от вас: %s.", effects_5_8[idx])
	}},
	{Range: DiceRange{Min: 9, Max: 12}, Value: func() string {
		return "В течение следующей минуты вы восстанавливаете по 5 Хитов в начале каждого своего хода."
	}},
	{Range: DiceRange{Min: 13, Max: 16}, Value: func() string {
		return "Существа совершат с Помехой спасброски от следующего требующего спасброска заклинания, которое вы сотворите в течение следующей минуты"
	}},
	{Range: DiceRange{Min: 17, Max: 20}, Value: func() string {
		idx := rander.Intn(len(effects_17_20))
		return fmt.Sprintf("Вы подвергаетесь эффекту, который длится 1 минуту, если в описании не указано иное.\n%s", effects_17_20[idx])
	}},
	{Range: DiceRange{Min: 21, Max: 24}, Value: func() string {
		return "В течение следующей минуты все ваши заклинания со временем сотворения в действие считаются имеющими время сотворения в Бонусное действие."
	}},
	{Range: DiceRange{Min: 25, Max: 28}, Value: func() string {
		return "Вы перемещаетесь на Астральный план до конца своего следующего хода; после этого вы вернётесь в ранее занимаемое пространство, или ближайшее свободное пространство, если это пространство занято."
	}},
	{Range: DiceRange{Min: 29, Max: 32}, Value: func() string {
		return "Не совершайте бросков урона для следующего наложенного вами в течение 1 минуты заклинания, наносящего урон. Вместо этого используйте максимальное значение для каждой кости."
	}},
	{Range: DiceRange{Min: 33, Max: 36}, Value: func() string {
		return "В течение следующей минуты у вас есть Сопротивление всему урону."
	}},
	{Range: DiceRange{Min: 37, Max: 40}, Value: func() string {
		return "Вы превращаетесь в растение в горшке до начала вашего следующего хода. На это время вы получаете состояние Недееспособный и Уязвимость всему урону. Если ваши Хиты опускаются до 0, ваш горшок разбивается, и вы возвращаетесь в свой облик."
	}},
	{Range: DiceRange{Min: 41, Max: 44}, Value: func() string {
		return "В каждый свой ход в течение следующей минуты вы можете Бонусным действием телепортироваться на расстояние до 20 футов."
	}},
	{Range: DiceRange{Min: 45, Max: 48}, Value: func() string {
		return "Вы и до трёх существ в пределах 30 футов от вас на ваш выбор получают состояние Невидимый на 1 минуту. Это состояние заканчивается для существа, как только существо совершает атаку или сотворяет заклинание."
	}},
	{Range: DiceRange{Min: 49, Max: 52}, Value: func() string {
		return "В течение следующей минуты рядом с вами парит призрачный щит, даруя вам бонус +2 к КЗ и иммунитет к Волшебным стрелам."
	}},
	{Range: DiceRange{Min: 53, Max: 56}, Value: func() string {
		return "В этот ход вы можете совершить одно дополнительное действие."
	}},
	{Range: DiceRange{Min: 57, Max: 60}, Value: func() string {
		idx := rander.Intn(len(effects_57_60))
		return fmt.Sprintf("Вы сотворяете случайное заклинание. Это заклинание не потребует Концентрации и действует свою полную длительность. %v", effects_57_60[idx])
	}},
	{Range: DiceRange{Min: 61, Max: 64}, Value: func() string {
		roll := rollDice(1, 4)
		return fmt.Sprintf("В течение следующей минуты любой воспламеняемый немагический объект, к которому вы прикасаетесь и который не несёт и не носит другое существо, получает 1к4(%v) урона Огнём и загорается.", roll)
	}},
	{Range: DiceRange{Min: 65, Max: 68}, Value: func() string {
		return "Если вы умрёте в течение следующего часа, то мгновенно вернётесь к жизни, как если бы на вас сотворили Реинкарнацию."
	}},
	{Range: DiceRange{Min: 69, Max: 72}, Value: func() string {
		return "Вы получаете состояние Испуганный до конца вашего следующего хода. Мастер определяет источник вашего страха."
	}},
	{Range: DiceRange{Min: 73, Max: 76}, Value: func() string {
		return "Вы телепортируетесь в видимое вами незанятое пространство в пределах 60 футов от вас."
	}},
	{Range: DiceRange{Min: 77, Max: 80}, Value: func() string {
		roll := rollDice(1, 4)
		return fmt.Sprintf("Случайное существо в пределах 60 футов от вас получает состояние Отравленный на 1к4(%v) часов.", roll)
	}},
	{Range: DiceRange{Min: 81, Max: 84}, Value: func() string {
		return "В течение следующей минуты вы излучаете Яркий свет в радиусе 30 футов. Все существа, заканчивающие свой ход в пределах 5 футов от вас, получают состояние Ослеплённый до конца своего следующего хода."
	}},
	{Range: DiceRange{Min: 85, Max: 88}, Value: func() string {
		roll := rollDice(1, 10)
		return fmt.Sprintf("Выберите до трёх существ, видимых вами в пределах 30 футов от вас. Каждое из существ получает 1к10(%v) Некротического урона, а вы восстанавливаете Хиты в количестве, равном сумме этого урона.", roll)
	}},
	{Range: DiceRange{Min: 89, Max: 92}, Value: func() string {
		roll := rollDice(4, 10)
		return fmt.Sprintf("Выберите до трёх существ, видимых вами в пределах 30 футов от вас; каждое из них получает 4к10(%v) урона Электричеством.", roll)
	}},
	{Range: DiceRange{Min: 93, Max: 96}, Value: func() string {
		return "Вы и все существа в пределах 30 футов от вас получаете Уязвимость к Колющему урону на 1 минуту."
	}},
	{Range: DiceRange{Min: 97, Max: 100}, Value: func() string {
		idx := rander.Intn(len(effects_97_100))
		return effects_97_100[idx]()
	}},
}

var effects_5_8 = []string{
	"[Монодрон](https://next.dnd.su/bestiary/21457-modron-monodrone/)",
	"[Дуодрон](https://next.dnd.su/bestiary/21456-modron-duodrone/)",
	"[Фламф](https://next.dnd.su/bestiary/21309-flumph/)",
	"[Единорог](https://next.dnd.su/bestiary/21601-unicorn/)",
}

var effects_17_20 = []string{
	"вас окружает тихая неземная музыка, её слышите только вы и существа в пределах 5 футов от вас",
	"ваш размер увеличивается на одну категорию",
	"у вас отрастает длинная борода из перьев, остающаяся на лице, пока вы не чихнёте — в этот момент перья разлетаются с вашего лица и исчезают",
	"вы должны кричать, когда говорите",
	"иллюзорные бабочки порхают в воздухе в пределах 10 футов от вас",
	"у вас на лбу появляется глаз, дающий вам Преимущество в проверках Мудрости (Восприятие)",
	"розовые пузырьки вылетают у вас изо рта, когда вы говорите",
	"ваша кожа приобретает ярко-голубой оттенок на 24 часа или пока эффект не будет снят заклинанием Снятие проклятия",
}

var effects_57_60 = []string{
	"[Смятение](https://next.dnd.su/spells/10203-confusion/)",
	"[Огненный шар](https://next.dnd.su/spells/10514-fireball/)",
	"[Туманное облако](https://next.dnd.su/spells/10283-fog-cloud/)",
	"[Полёт](https://next.dnd.su/spells/10518-fly/)",
	"[Намасливание](https://next.dnd.su/spells/10532-grease/)",
	"[Левитация](https://next.dnd.su/spells/10560-levitate/)",
	"[Волшебная стрела](https://next.dnd.su/spells/10567-magic-missile/)",
	"[Отражения](https://next.dnd.su/spells/10584-mirror-image/)",
	"[Превращение](https://next.dnd.su/spells/10604-polymorph/)",
	"[Видение невидимого](https://next.dnd.su/spells/10632-see-invisibility/)",
}

var effects_97_100 = []func() string{
	func() string {
		roll := rollDice(2, 10)
		return fmt.Sprintf("вы восстанавливаете 2к10(%v) Хитов.", roll)
	},
	func() string {
		roll := rollDice(2, 10)
		return fmt.Sprintf("один союзник восстанавливает 2к10(%v) Хитов.", roll)
	},
	func() string {
		return "вы восстанавливаете ячейку заклинаний."
	},
	func() string {
		return "союзник восстанавливает ячейку заклинаний."
	},
	func() string {
		return "вы восстанавливаете все Очки чародейства"
	},
	func() string {
		return strings.Join(effects_17_20, "; ")
	},
}
