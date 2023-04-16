<p align="left">
    <img property="og:image" src="https://repository-images.githubusercontent.com/577755312/57f67b11-437b-448f-b53e-cf47165612c2" width="25%">
</p>

# Chicha 2 - это хронограф, используемый для соревнований.

<img align="left" property="og:image" src="https://upload.wikimedia.org/wikipedia/commons/thumb/b/b4/2012%2C_Casey_Stoner.JPG/548px-2012%2C_Casey_Stoner.JPG" width="15%"> 

"Эй, ребята! У нас есть крутой хронограф, который специально разработан для проведения соревнований во всех видах спорта: от бега и плавания до авто-мото-спорта и морских соревнований. Он может обслуживать более 100 000 участников одновременно и даже совместим с технологией UHF-RFID! И самое лучшее - мы готовы предложить его вам абсолютно бесплатно!"

> Обратите внимание, что при использовании технологии UHF-RFID для точного подсчета результатов быстро движущихся участников необходимо замедлить их скорость до 25-30 км/ч в зонах замера. Для этого рекомендуется установить специальные считывающие устройства (финишные или промежуточные ворота) на крутых поворотах, где скорость движения не превышает 30 км/ч. Такой подход гарантирует точность считывания RFID меток.

<p align="left">
    <img property="og:image" src="https://repository-images.githubusercontent.com/368199185/e26c553e-b23e-4bae-b4d2-c2df502e9f04" width="75%">
</p>


### Демо версия: [http://chicha.zabiyaka.net](http://chicha.zabiyaka.net/)

- Скачайте последнюю версию [↓ CHICHA.](http://files.matveynator.ru/chicha2/latest/)
- Для запуска гонки в тестовом режиме - скачайте [↓ RACETEST.](http://files.matveynator.ru/racetest/latest/)

> Поддерживаемые операционные системы: [Linix](http://files.matveynator.ru/chicha2/latest/linux), [Windows](http://files.matveynator.ru/chicha2/latest/windows), [Android](http://files.matveynator.ru/chicha2/latest/android), [Mac](http://files.matveynator.ru/chicha2/latest/mac), [IOS](http://files.matveynator.ru/chicha2/latest/ios), [FreeBSD](http://files.matveynator.ru/chicha2/latest/freebsd), [DragonflyBSD](http://files.matveynator.ru/chicha2/latest/dragonfly), [OpenBSD](http://files.matveynator.ru/chicha2/latest/openbsd), [NetBSD](http://files.matveynator.ru/chicha2/latest/netbsd), [Plan9](http://files.matveynator.ru/chicha2/latest/plan9), [AIX](http://files.matveynator.ru/chicha2/latest/aix), [Solaris](http://files.matveynator.ru/chicha2/latest/solaris), [Illumos](http://files.matveynator.ru/chicha2/latest/illumos)

- Download latest version of [↓ CHICHA.](http://files.matveynator.ru/chicha2/latest/) 
- For race testing - download [↓ RACETEST.](http://files.matveynator.ru/racetest/latest/) 

### Хронограф может быть использован в двух режимах: "mass-start" и "delayed-start".
> Существует опция конфигурации под названием "-race-type", которая позволяет выбирать тип гонки: "mass-start" (масс-старт) или "delayed-start" (отложенный старт).

Режим "mass-start", - словно сигнал трубы: все участники сразу в бой! А режим "delayed-start" дает возможность стартовать поочередно, с небольшим перерывом между началом гонки. Например, в забеге на 100 метров или в мотокроссе, спортсмены могут стартовать все вместе (mass-start) или же в авторалли последовательно, с некоторым временным интервалом (delayed-start). Исходя из условий проведения соревнований, выбирается соответствующий режим работы хронографа.


### ![#FF0000](https://via.placeholder.com/15/FF0000/000000?text=+) ![#008000](https://via.placeholder.com/15/008000/000000?text=+) ![#EE82EE](https://via.placeholder.com/15/EE82EE/000000?text=+)  Цветовые подсказки во время гонки:

> В авто и мотоспорте, на соревнованиях, на которых спортсмены борются за лучшее время круга или за наилучший результат в гонке, используется система цветовых сигналов на табло для показа изменений времени круга. Когда спортсмен завершает круг, его время отображается на табло, и цвет сигнала указывает на то, улучшил ли он свой результат по сравнению с предыдущим кругом или нет. 

Вот как работает алгоритм:

![#008000](https://via.placeholder.com/15/008000/000000?text=+) Зеленый цвет: если время круга лучше предыдущего, то на табло будет отображаться зеленый цвет. Это означает, что спортсмен улучшил свой результат, и это может стимулировать его на дальнейшее улучшение времени.

![#FF0000](https://via.placeholder.com/15/FF0000/000000?text=+) Красный цвет: если время круга хуже, чем предыдущее, на табло будет отображаться красный цвет. Это означает, что спортсмен ухудшил свой результат, и ему нужно работать над улучшением.

![#EE82EE](https://via.placeholder.com/15/EE82EE/000000?text=+) Фиолетовый цвет: если на табло появляется фиолетовый цвет, это означает, что спортсмен показал лучшее время круга на трассе. Это может быть достигнуто в конце сессии, когда все спортсмены завершают свои круги, или в середине сессии, если спортсмены уже успели улучшить свои результаты.

Цветовые сигналы на табло используются для помощи спортсмену в оценке своей производительности и понимании, насколько он улучшает свои результаты. Это также помогает зрителям понимать, как проходит гонка и кто лидирует.


### Вспомогательные конфигурационные опции:
```
chicha -h
Usage of chicha:
  -average
    	Calculate average results instead of minimal results.
  -average-duration duration
    	Duration to calculate average results. Results passed to reader during this duration will be calculated as average result. (default 1s)
  -collector string
    	Provide IP address and port to collect and parse data from RFID and timing readers. (default "0.0.0.0:4000")
  -db-path string
    	Provide path to writable directory to store database data. (default ".")
  -db-save-interval duration
    	Duration to save data from memory to database (disk). Setting duration too low may cause unpredictable performance results. (default 30s)
  -db-type string
    	Select db type: sqlite / genji / postgres (default "genji")
  -lap-time duration
    	Minimal lap time duration. Results smaller than this duration would be considered wrong. (default 45s)
  -pg-db-name string
    	PostgreSQL DB name. (default "chicha")
  -pg-host string
    	PostgreSQL DB host. (default "127.0.0.1")
  -pg-pass string
    	PostgreSQL DB password.
  -pg-port int
    	PostgreSQL DB port. (default 5432)
  -pg-ssl string
    	disable / allow / prefer / require / verify-ca / verify-full - PostgreSQL ssl modes: https://www.postgresql.org/docs/current/libpq-ssl.html (default "prefer")
  -pg-user string
    	PostgreSQL DB user. (default "postgres")
  -proxy string
    	Proxy incoming data to another chicha collector. For example: -proxy '10.9.8.7:4000'.
  -race-type string
    	Valid race calculation variants are: 'delayed-start' or 'mass-start'. 1. 'mass-start': start time is not taken into account as everybody starts at the same time, the first gate passage is equal to the short lap, positions are counted based on the minimum time to complete maximum number of laps/stages/gates including the short lap. 2. 'delayed-start': start time is taken into account as everyone starts with some time delay, the first gate passage (short lap) is equal to the start time, positions are counted based on the minimum time to complete maximum number of laps/stages/gates excluding short lap. (default "mass-start")
  -timeout duration
    	Set race timeout duration. After this time if nobody passes the finish line the race will be stopped. Valid time units are: 's' (second), 'm' (minute), 'h' (hour). (default 2m0s)
  -timezone string
    	Set race timezone. Example: Europe/Paris, Africa/Dakar, UTC, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones (default "UTC")
  -version
    	Output version information
  -web string
    	Provide IP address and port to listen for HTTP connections from clients. (default "0.0.0.0:80")
```
