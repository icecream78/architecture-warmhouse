# Project_template файл для проектной работы 1-го спринта

# Задание 1. Анализ и планирование

«Тёплый дом» — это небольшая компания, которая организует удалённое управление отоплением в доме. Недавно она выиграла тендер и получила заказ на создание экосистемы умных посёлков на территории нескольких регионов страны.

Для выполнения заказа компании нужно расширить функционал своей системы и масштабировать её. Пользователи должны иметь возможность управлять не только отоплением, но и умными устройствами в доме. Чтобы подключиться к экосистеме, в домах должны быть установлены специальные датчики и реле. Сейчас они есть только в половине домов, которые хотят подключиться.

Состояние компании в настоящий момент не позволяет в полной мере реализовать новые бизнес-цели. Для этого требуется пересмотр и оптимизация всей экосистемы.

### 1. Описание функциональности монолитного приложения

Основная деятельность пользователя автоматизируется через веб-приложение, доступное через сеть интернет.

**Управление отоплением:**

- Пользователь может регулировать температуру внутри дома (начиная от граничных случаев - включить/выключить отопление, так и установку промежуточных значений) с точностью до комнаты;

**Мониторинг температуры:**
- Пользователь может получить данные о температуре внутри дома, с точностью до определенного комнаты/конкретного датчика. В случае временной/полной недоступности датчика, пользователь увидит устаревшие данные с последней датой снятия показаний;
- Пользователь может получить данные о температуре в разных единицах измеренция - по цельсию/фаренгейту;

**Обслуживание системы управленния:**
- Пользователь может подключить новые датчики только путем вызова специалиста, который может создать новый датчик/отключить неисправный;

### 2. Анализ архитектуры монолитного приложения

Язык реализации: Golang
Для хранения данных используется PostgreSQL
Архитектура системы: единое монолитное приложение
Взаимодейтвие: синхронные http запросы без намёков на параллелизацию, но с идеей graceful degradation в случае неуспешных запросов
Масштабируемость: масштабировать можно как горизонтально (доп. ресурсами), так и вертикально (количеством реплик), но в таком случае утилизация ресурсов будет неравномерной из-за разного паттерна нагрузки на части системы
Развертывание: обновление системы происходит за счет полной остановки старой версии и запуска новой

Список текущих особенностей и проблем системы:
1. Система имеет единый init.sql, в котором описывается схема данных всего приложения. На текущий момент нет никакой автоматизации в виде миграторов для обеспечения предсказуемого процесса изменения схемы данных, а также истории почему и как схема менялась
2. Система не имеет четко описанного API для взаимодействия с ней. Из-за этого не ясен итоговый список доступных команд, их контрактов и функциональности
3. Слой handler включает в себя как знание о текущем http фреймворке, так и содержит части бизнес-логики
3. Смешение логики работы моделей системы с особенностями API и тем, как данные расположены в БД. Бизнес-логика должна быть максимально независимой и не должна что-то знать о деталях реализации
4. При получении информации по датчикам, происходят последовательные синхронные запросы за информацией о температуре без намеков на параллелизацию. В добавок эти походы ограничены клиентов в 10 секунд, что в итоге может привести к росту времени ответа всей ручки и увеличивает риски, что пользователь попросту не долждется результата ответа
5. Некоторые конфигурационные параметры унесены в переменные окружения, при этом другая их часть остается захардкорженной, что делает систему менее гибкой и предсказуемой
6. В системе отсутствует логгирование, отправка метрик, оповещение об ошибках работы система. Из-за этого не всегда понятно работает ли система или нет. Не спасает даже /health метод, ведь он не проверяет работу важного бизнес-функционала или же, хотя бы, доступность основных компонентов системы
7. Не совсем очевидное поведение в лице PUT /sensors/:id, который бесконтрольно может обновить всю информацию о датчике
8. В системе отсутствует система валидации входящих данных, что может привести к записи невалидного контента в бд, либо же к выходу из строя конкретных датчиков, которые например, не поддерживают возможность выставить температуру выше 70 градусов по цельсию, либо же датчик работает по шкале фаренгейтов и пользователь задаст не совсем корректное значение
9. Система не имеет никакой авторизации, acl и прочих вещей, которые позволят защитить систему от несанкционированного доступа

### 3. Определение доменов и границы контекстов

Основные домены As-Is решения.
*Подключение и обслуживание датчиков*
- поддомен подключение/отключение новых датчиков
1. контекст: подключение нового датчика
2. контекст: отключение существующего датчика
- онбновлений данных о датчиках
1. калибровка показаний по датчику
2. обновление данных по датчику

*Управление отоплением в доме*
- поддомен установка температуры
1. контекст: установка температуры в доме/комнате

*Мониторинг показаний датчиков*
- поддомен отслеживания показаний температуры
1. контекст: получение информации о температуре по датчикам

### **4. Проблемы монолитного решения**

- нет четкой документации API
- слабые границы бизнес-логики приложения и наличие обобщенного кода
- смешение логики хранения данных и взаимодействия с этими данными
- недостаточные возможности масштабирования приложения и отсутствие возможности масштабировать конкретные части приложения
- невозможность обновления системы без простоя
- отсутствие логики лимитирования времени обработки внешнего запроса

Резюмируя, в будущем у приложения будут как проблемы с производительностью, так и возникнут сложности с поддержкой и добавлением нового функционала.

### 5. Визуализация контекста системы — диаграмма С4

С диаграммой контекста в нотации C4 можно ознакомиться по [ссылке](https://www.planttext.com?text=VLDDRzD04BtxLomv1Iaq5qwSAWLkWAeWn35oaoMnP7iZUoF2hMduL55oGO9BXO0_u4BI9h6J_8NPVyIRjLDAi7hmship-zwRD-F3MBkH7WI-eptjIr6X7pIdBQNcZ9Q2PhIXh28QAjHogCL3p-r6Rk0uMlM5Lk9OQQq2qV4YsTQU2XtdZXUs_K573Y9VzBJknm_gzSXzTT3rT6zmF8Xbr6QiK1-qLL1lUcqtwgYddknBlzwuV-_8TqFz_CdZUk36MaOuQZLKq5SXt-YpJET8Hh4AgmnLWznlK9YQJtI5zozAp2daOr-v9IRCY1PcPBUPEHClt2ZeAG1MxkxGNth37FwS4ahnpCnxMj3Amdw7FtFp3dkEMII1eurTODaA91FapscDw6IVZZRutB37HSOKlCzUezwWUYqbtzVA-y77NsbjtW3gX39Knc8s05lZU_eF8IKWf0iSfMPCbyXZ9oRceyqiKRm6qxnWuijRrWkQELq0rLg8PmdDXUZMew6knZACbNNiCnkTkqMBic9BajWygDk29e6-zDewbiv58ZZBcyW1heiLfbo1pQ3VLY7YhCjfIhyOa193V-xpSBaxcOlTa5WVibi4W8XfyadPfybdxw-n2GRcb2NVABiciDe9GdkOmAOipDkyAnoxokcuJBZ4pHtPNcneMI1mHiPbc6a9XRt32vHpMzaSY6urlMvF_pTiHGyhz0MWrYFp8PyN2ssTMPanHiTQLYBo9da2TGbVIh2vn0NS7DdjkZ1_3pgplzM9XHYp8L1xACCE_le_)

# Задание 2. Проектирование микросервисной архитектуры

В этом задании вам нужно предоставить только диаграммы в модели C4. Мы не просим вас отдельно описывать получившиеся микросервисы и то, как вы определили взаимодействия между компонентами To-Be системы. Если вы правильно подготовите диаграммы C4, они и так это покажут.

**Диаграмма контейнеров (Containers)**

С диаграммной контейнеров To Be можно ознакомиться по [ссылке](https://www.planttext.com?text=tLXBJnjN5DxxLzooAKY0DrrrbGIIz95A3b1bYSRs14QrPwpnE2cg8eMu944WQDD5DsdIBBgtZLqcP73-mim_wdVElVEUCMuKLL65CFW-plCxhtkjwsYsqsktAfyPPgFLQ-fYrt4wtI-hLLjxkBfZEBkzUg-hsmtBT7JJMMrOxMgdfPbqPwNhD5j6lRhsnKgDbhx_xawrhNKTtSRIzXfkQ8QfswiTeb_vNFY_1O-3llyUFox_D_4zWIzJ_rrmx9_ZxyWV14-m_0w_RmM-UluasASb4UpBevu_3Gw3u-14tB1wPbEpzohXFuiaw6tRwdMMnOB_kf2k2qBWpZHzBpXQM18_LOIewNRNCfb0wTtWE7DNtDFizgw5ImiVT6afPgmrsuP9HDv0oKdmr7Vz0KbBzwJUO_nCWiV-t_w5xowK70Nv_E5bWSyf-4zXz2EOChIZVmOfJy20R3lrXq4_U0OXBq6HFND1P6c1xleZ0UCH6_yzA0-WuhlaOLe43LvKNt7y4Xn8Bjfx2uAJdFahZnWQh_Z6XEI6qKubaJDc11R1qw0FEE0kuI9DehAv1yYrjoCKH0QM6z9KqGL_aFG0Ery82F07E0t9H4fsA3A5_QI6RqjW6foA5S52dZDK3_4VIpuYy-EC2vq7VE_IdpAhA2osKx2zfJNqvlLEfqlYpfO5_XnZaynAmNGP70K7PVATH2ecE2y-rEjWHPnUuDImYCj4GZg_lgpqvSyXuufS1FK8MMEsIRn7W1onEVpX2X263GXzpk2Sb1kLe4DxvrYb4wU0wI7s9WI28taDYf-20UAhdqTa21Olf3GcD8EA13Xinbj0V8HE1WY4e0rI_XljWRRPi8sEiomsTAtXSBoc3TcswaPBLxPyoVXtGs-pVddxYKNDRDgMqQmQLhUgEusbZsTfH9iWPF2PKOWCBocNN7hblraLdB0Y2yCbJqfHz1-Q_LQhvw22A6C_uizvOtG0VY336KCo6CYqvm_AOeyJhtBWXuITT0CRkyySMKzydlbJbItJ9cNXJWl3acMvDqTi2ddq8o9dHeogu_lvQYEBGc5nAWkUj3zVKdua611ey7V4T8_AV5NN6Z_ePdDJjnyO3JsRARqYObB42wvtxwNn8M2HyUiQMfY6BF-_ieuupPM6O3f2-W5vmHIfIhdiP0O6U-5FAdJP6aOHXPYaEiOMlQ0gHa6Qgc1zcUxKLKQOZAgG11R3sCCrrLPmC570SH1xxFvBTPjTiAOrctlLRorp9qfSEQLDop7kxoc3Vg3clz1HN3ZWun7AfCWZBefHTYZjzLXWz9UsSR-tO-KpxNfToRfU9vv_HOq1csUWSWknBSjXnn0Ju0nBvPYRK5eQ4E-QrNLsR7tpxXs08p9KWSaWI5CdQwtVk5AGCIVUE46gb95fYn8xOv8MkPAxXk0WArY95o6IrdDsvn2efDoMSChdfuuHG-LwxMjijJjMrt0ConJ9g6DSl81EOR8SI22F16jCNNYZv2fLC9LJfp8NgGHsGNOR1GVaEuXyoAxyQckh9g3EGBQVCjw9-sqBqJm3-6Im1EpVn09G5Hcfgk4IQSvvrFaAZj2pH6rHEoLcPIkyukesWLEL3Suf_0jBPmaPMKn7gPmXsxEhb9H0HMUfDVDQWbkY1b1Zdlgodq5ZM9x_ZK8v0gn6NyK5vGbFRDJMHF0pMzkJEi-XsTsUtjCNzGT863TwNQdQXSg455eeLyyvJmqICG2KXpL8ZXmlA97Xe9IsUAYu-QEps6rZH9QsJ9SSgEAfqCrtD0lfiuYT5GdFQAPrYH7VdK5cEOw3S-oUAE27TwZh8NDwq8rr0Xdf1z8gQZSgTxIzxNKrtihUWPSsaYDpQWnPIbn9p2PXbPJS20jKSvvBdrFGGZaKeZ4LjPbuTWiQ6xbNsT1RIjBad491xdATe_Agh4Dl23BHSO-IGo59D5R3hYUAcr8rltBBNE846IqRDpUt1D8UJ1sBc9e0Fa4frSD2hhMS9IZXDDFmgP4crMDJLqLTSNvGS-E1b9uKN9uNajqpMZ9yvwxeZCJ8vOh82lJSyunk9nuZ02xAJzctXQaQqvCcc0EbnIXa0OPbtH09nkOD4Phi7Zz152HG-ebBcWBg_BB-xo2RcE5bSHe1b-TSGgDHC_kMoGQCRPPgr_pVuoQ9NmbLUqFDXBJ14soSoZR_Z50lHH2HyrqE1I_qKu9uji_8iYBNv45q8kBUANfuv5Ke67S4J-N4bcEQwU0oF7d8YYVSVlHyN3HzUFCaYinKwHK2cZt2hSxCTmvM9csrJ9DnHMGoP6fh97UUL78-AQ4JDPn2jljitAH1uu5YJWfv4idsx3LxWThdOJpMSWlGvoTZcjpxuRpC6KE2dbgV8sgLwCMaGDMFpuEPPCRYvEHvfJcsrcf4vffkDdljrZy0)

**Диаграмма компонентов (Components)**

С диаграммой внутреннего устройства компонента по работе с устройствами умного дома можно ознакомиться по [ссылке](https://www.planttext.com?text=rLTDRzj64BthLsnzYGN44IBar5CI9Ma2980I9sU3fZQs08fI-P6sDykkwnGA8Y2N007g9eVUPSAAQLkc_yBkFygR8McHD1cxHT2wFj6xipFllTcPMTrsFSFn_9ubMhTMxbko9stFhNtNjKtBxqYnxda3z_j6mp5-kh7MzTRzLT-LZjctFTZTCFkzni0oR3g_v7eTgxlQQDrQQjFIYoUF6pt3zQI3fPLMlpVestHc0DjQJOWMV1XTMpgBng2xB9rNNLDU5mjtsey46xjThzktXHwggTvGWGhrS07xwfqUwQ7Ur5lgL4smjobK860HgMCL2hFYg9e8jQ-EzLZ0S0-h4xMd8hsfZd181NeRo_jwh7-7yRPu-FnvMwW369sIPnoFL43HV-Y3sjf2FIQGK5hiI0B_m32D3b5Gd-Pn1L168XDlMlXdIBWcuBEZGe2Pg9aAcUSkzWCusabP0euWrdm-0DkfEi6PQOa3bcGBw-0AeqZDnBDsI-XV4NK31XEAdUCZH0cZXuRTiQJpeEyyCMp_fM5wlYDlVWk612KaA7eaRluxTCKYlf1a-3dMRp8O0I6RwWRZ9YYRzO94bzAd-QtrQLu1VKv96LX5UWm7v8BZ1F1waaWt1n5lSz5COHHY8S5MhwoXfxxqvQ9yXS-xlll8Hb_WObSdUWiGDcB7AD2nqA_P9HVqX5T7ADdgwkR31q0AMbWYib7SCZ84ivhnRbR_kk2_StXNxl_idM76Tvpv8mWH8Js0M8KsVvlkhIRtwTvTyl0N8uDUbEMpGoHu998SPhIEoIF556x8KaNPdTnbG6H4QFTTRysHoqyVKnAk2L8PT-frSd-Gk92JFINwjULVqCDxJDNjOLJ4KBaB_94ek46tQvP4HBvI68ViXXjdnXRjVYdkg3WpPPEG8l7PhxYvFaykigdV9hqwFQCEOtLc-TvIho3G_7yGQ2OfsaK9KbP623gXp52Q0bJZ0UKUCDHdpA-jmltFjvHpWEZUq6YaxGgJUBFsJ5gvCLsSUE_d_JDf5oMzIR_bYtfAengHTZAhu3zlHh_X_zZkN2-ZJ3W34w8GeZoAYsPM0IFNP04B5LaP5lJeCccO3hMPbEWHPO0Qo8pBThzCXF4_vD2ycXmo9AgaBxQDgqJ2xmnoYQYIljWvFinzWqELTyzClUVFNvJcZyM9HJe5A7wMAPsdDEwo7EEYcVAyeq4UafgeRmn1z9QNbk-knq9JAHt4CqOTpd6Mznbya6ruzZf8LUTPmTCj6Q3PSPHxUzD0nMSwoiDBt9bpBu7TR7RfwJFaLZVYvdkOT0jgUKDEmMUi9Wp0PLRHsV8v-dfM_Wj4XPZ5hai_XNRedS3nkGOAhvzoL1U5hcnCyGFbpxDQJPwYIFTszQD5xLNAK3ecqsUA-wDLZgNmuhfC7d8lxh-59gFARMbtyEFr2m00)

**Диаграмма кода (Code)**

Добавьте одну диаграмму или несколько.

С диаграммой кода компонента по работе с температурой можно ознакомиться по [ссылке](https://www.planttext.com?text=rLJTIiCm5BxFK-JE3PqFi8injZ4Hl7A-G5WFCv0s8ys48S9-BaGc23x9A8vDrFaA9s_aoTgwLGUYTr8AwPVl-yxNqvpMOokC7OIAMMaLyAOE-pg2o5OvFkARcxg96-42t-XQu88JCDv0QtpsAqnmpTqO5mGbkABLdECxORYcpIUSkrklcf3ubN1FcvD-waPknbXNYJZc7OXYRVWruxp2PL3b7MjarDiysc4V2e3INM6bZWfW077No7u1lydjBqKKA4Zz3_a9s0i8-s247HXe6nsUPOwbt3fpA5TvQOVQ4bOqXcdM8mx6Q5F-oL46l_2Ym5dY6CoLxC9_I9m63i3dRIf91-TFgSz1A2l3JPcsDZlAnZ3NVIyRMsr4RwTFQRuMM25Ll5NkgnZg0DJn_kPyiJbj5TkoLqgX3lmW9Eu1bxXY-QchL6gVtvDcAd4pNDATu9o6OS9hdEMx6sRsHiJqiqI3X2_kdXJJb5qW5xg-pvkntC_Jip7pHB9T4jtFujsMpEeG1VGh-000)


# Задание 3. Разработка ER-диаграммы

С ER-диаграммой можно ознакомиться по [ссылке](https://www.planttext.com?text=VL91QiCm4Bph5NlhoOTSJGW9cD92yINE0nHxwwYr9IF9eQRaLRsKakODUbBNSXhB8UHWyCwEPcPNMXhBjSrtA1b3QONjO6DGmoS3U4vY4D9YIVLy_exTOa5eoclqRO17IVyn6Ak5B3toSeKSw5ljENa4u1e_OjWgLVND4Yycx739yAHQWtT2h8f2ep61wAeIFAmJDBaMZHLA1YXDQzkHDXeck1Vvw3Zq0yCrQi6hjAst68xYSmOH2Sgw9jp0HWeRbBcIhw9iDU-3JGzUPQDbgfnFybDSfh7oeDc9vfnwxGUziws1a67Tq5aAzsQKlBcCYNY65TPeTPVG_NdrmMxSJmzHxoAOYHf9j6vYACYeLtoOWVjn9t17z-jExijzkoFpBonAe_CrtoODlPurs0uicB7pLKMpQ3B_RPL_CibCvc5iYRBmx_uF)

# Задание 4. Создание и документирование API

### 1. Тип API

Укажите, какой тип API вы будете использовать для взаимодействия микросервисов. Объясните своё решение.
Все взаимодействие в рамках системы можно разделить на 2 вида:
- доступ к системе извне
- доступ внутри системы
Для работы с системой извне используется REST API, т.к. он отлично себя зарекомендовал при создании CRUD приложений, которым отчасти является и наше. Это API будут использовать клиентские приложения.
Для синхронного взаимодействия внутри системы будет использоваться RPC API на базе protobuf для обеспечения высокой производительности и упрощенной интеграции разных компонентов между собой.
Для асинхронного взаимодействия будет использоваться также protobuf из-за своей эффективности и простоты.

### 2. Документация API

Здесь приложите ссылки на документацию API для микросервисов, которые вы спроектировали в первой части проектной работы. Для документирования используйте Swagger/OpenAPI или AsyncAPI.
Всего было описано 4 метода из общего ландшафта системы - 2 для внешнего потребления и 2 для внутреннего.
Для внешнего взаимодействия с ситемой вызывается REST API, схема которого доступна в виде файла [rest.openapi.yml](schemas/rest.openapi.yml)
Для внутреннего потребления используется protobuf и [здесь](schemas/session.proto) можно ознакомиться с примером описания сервиса, а [здесь](schemas/device_disconnected_event.proto) с примером события, которое будет отправляться в общую шину данных.

# Задание 5. Работа с docker и docker-compose

Перейдите в apps.

Там находится приложение-монолит для работы с датчиками температуры. В README.md описано как запустить решение.

Вам нужно:

1) сделать простое приложение temperature-api на любом удобном для вас языке программирования, которое при запросе /temperature?location= будет отдавать рандомное значение температуры.

Locations - название комнаты, sensorId - идентификатор названия комнаты

```
	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}
```

2) Приложение следует упаковать в Docker и добавить в docker-compose. Порт по умолчанию должен быть 8081

3) Кроме того для smart_home приложения требуется база данных - добавьте в docker-compose файл настройки для запуска postgres с указанием скрипта инициализации ./smart_home/init.sql

Для проверки можно использовать Postman коллекцию smarthome-api.postman_collection.json и вызвать:

- Create Sensor
- Get All Sensors

Должно при каждом вызове отображаться разное значение температуры

Ревьюер будет проверять точно так же.


