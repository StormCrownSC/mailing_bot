package cconstant

var (
	RegisterStep = map[uint8]string{
		0: "Введите telegram id: ",
		1: "Введите имя: ",
	}
	HelpText = "## Команды бота\n\n- **/sending [textID]**  \nОтправить рассылку с текстом по указанному \"text id\"\". Если текст с таким ID не найден, бот сообщит об ошибке.\n\n- **/adduser **\nНачать процесс добавления нового пользователя. Введите данные пользователя по запросу бота. Чтобы отменить добавление, введите \"stop\"\".\n\n- **/removeuser [telegramID]**\nУдалить пользователя с указанным \"telegramID\"\".\n\n- **/recoveruser [telegramID]**\nВосстановить пользователя с указанным \"telegramID\"\".\n\n- **/getusers**\nПолучить список всех пользователей.\n\n- **/addadmin**\nНачать процесс добавления нового администратора. Введите данные администратора по запросу бота. Чтобы отменить регистрацию, введите \"stop\"\".\n\n- **/removeadmin [telegramID]**\nУдалить администратора с указанным \"telegramID\"\".\n\n- **/recoveradmin [telegramID]**\nВосстановить администратора с указанным \"telegramID\"\".\n\n- **/getadmins**\nПолучить список всех администраторов.\n\n- **/addtext [text]**\nДобавить новый текст для рассылки. Убедитесь, что текст не слишком длинный.\n\n- **/removetext [id]**\nУдалить текст рассылки с указанным \"id\"\".\n\n- **/recovertext [id]**\nВосстановить текст рассылки с указанным \"id\"\".\n\n- **/gettexts**\nПолучить список всех текстов рассылок."
)
