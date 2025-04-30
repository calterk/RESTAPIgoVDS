function Main() {
    Util.call(this); // Наследование от Util
    const util = new Util();

    this.init = () => {
        console.log("init");
        document.querySelectorAll(".action").forEach(el => {
            console.log(el);
            el.onclick = () => {
                this[el.dataset.action]();
            };
        });
        this.checkTokenValidity();
    };

    this.checkTokenValidity = () => {
        const token = localStorage.getItem("token");
        // Отправляем запрос на сервер для проверки токена
        util.get("/validateToken?token=" + token, data => {
            console.log(data, "saflknadsfl");
            if (data === true) {
                // Токен валиден, отображаем таблицу пользователей
                this.showUsersTable();
            } else {
                // Токен невалиден, отображаем форму авторизации
                this.showAuthForm();
            }
        }, error => {
            console.error("Error validating token:", error);
        });
    };

    this.showAuthForm = () => {
        util.id("authForm").style.display = "block";
        util.id("usersTable").style.display = "none";
        util.id("vdsTable").style.display = "none";
        util.id("vdsTableAdmin").style.display = "none";
        util.id("addUserBtn").style.display = "none";
        util.id("deleteUserBtn").style.display = "none";
        util.id("addVdsBtn").style.display = "none";
        util.id("addVdsBtnAdmin").style.display = "none";
        util.id("deleteVdsBtn").style.display = "none";
        util.id("logoutBtn").style.display = "none";
    };

    this.showUsersTable = () => {
        util.id("authForm").style.display = "none";
        util.id("logoutBtn").style.display = "block";
        console.log("табла пользователей");
        this.getUsers(); // Получаем и отображаем пользователей
    };

    this.showVdsTableAdmin = () => {
        console.log("табла vds admin");
        this.getVdsAdmin();
    };

    this.showVdsTable = () => {
        console.log("табла vds");
        this.getVds();
    };

    this.addUser = () => {
        util.modal('addUserModal', 'show');

        document.getElementById('saveUserBtn').onclick = () => {
            const user = {
                login: util.id('addUserLogin').value,
                pass: util.id('addUserPassword').value,
                name: util.id('addUserName').value,
                access: util.id('addUserAccess').value === 'true'
            };

            util.post('/addUser', user, (response) => {
                if (response.success) {
                    this.getUsers();
                } else {
                    alert(response.message);
                }
            });

            util.modal('addUserModal', 'hide');
        };
    };

    this.addVds = () => {
        util.modal('addVdsModal', 'show');

        document.getElementById('saveVdsBtn').onclick = () => {
            const vds = {
                name: util.id('addVdsName').value,
                img: "//1111",
                ip: util.id('addVdsIp').value,
                login: util.id('addVdsLogin').value,
                pass: util.id('addVdsPassword').value,
                status: false
            };

            util.post('/addVds', vds, (response) => {
                if (response.success) {
                    this.getVds(); // Заменили на getVdsAdmin
                } else {
                    alert(response.message);
                }
            });

            util.modal('addVdsModal', 'hide');
        };
    };

    this.addVdsAdmin = () => {
        util.modal('addVdsModalAdmin', 'show');

        document.getElementById('saveVdsBtnAdmin').onclick = () => {
            const vds = { // Изменили user на vds
                Uid: parseInt(util.id('addVdsUidAdmin').value,10),
                name: util.id('addVdsNameAdmin').value,
                img: "//1111",
                ip: util.id('addVdsIpAdmin').value,
                login: util.id('addVdsLoginAdmin').value,
                pass: util.id('addVdsPasswordAdmin').value,
                status: false
            };

            util.post('/addVdsAdmin', vds, (response) => {
                if (response.success) {
                    this.getVdsAdmin();
                } else {
                    alert(response.message);
                }
            });

            util.modal('addVdsModalAdmin', 'hide');
        };
    };

    this.deleteUser = () => {
        // Вызываем модальное окно
        util.modal('confirmDeleteModal', 'show');

        // Получаем кнопку подтверждения удаления и добавляем обработчик событий
        document.getElementById('confirmDeleteBtn').onclick = () => {
            const userId = parseInt(document.getElementById('deleteUserIdInput').value, 10);

            util.post('/deleteUser', {uid: userId}, (response) => {

                // Обработка ответа сервера
                if (response.success) {
                    // Обновить таблицу пользователей
                    this.getUsers();
                } else {
                    // Ошибка при удалении пользователя
                    alert(response.message);
                }
            });

            // Закрыть модальное окно
            util.modal('confirmDeleteModal', 'hide');
        };
    };

    this.getUsers = () => {
        console.log("Data from /getUsersTable");
        util.get("/getUsersTable", data1 => {
            if (data1 === "f") {
                console.log("agaaaaaaa")
                this.showVdsTable();
            }else{
                console.log(data1, "Data from /getUsersTable");

                if (Array.isArray(data1)) {
                    util.id("usersTable").style.display = "block";
                    util.id("addUserBtn").style.display = "inline-block";
                    util.id("deleteUserBtn").style.display = "inline-block";
                    const tableBody = util.id("usersTable").querySelector("tbody");
                    tableBody.innerHTML = ""; // Очищаем текущие строки

                    data1.forEach(user => {
                        const row = `<tr>
                    <td>${user.uid}</td>
                    <td>${user.Login}</td>
                    <td>${user.Pass}</td>
                    <td>${user.Name}</td>
                    <td>${user.Access ? 'Да' : 'Нет'}</td>  
                </tr>`;
                        tableBody.innerHTML += row;
                    });
                } else {
                    console.error("Data is not an array:", data1);
                }
                console.log("blya", data1)
                this.showVdsTableAdmin();
            }
        }, error => {
            console.error("Error fetching data:", error);
        });
    };

    this.getVds = () => {
        util.get("/getVdsTable", data1 => {
            console.log(data1, "Data from /getVdsTable");
            if (Array.isArray(data1)) {
                util.id("addVdsBtn").style.display = "inline-block";
                //util.id("deleteVdsBtn").style.display = "inline-block";
                const cardsContainer = document.querySelector(".row");
                cardsContainer.innerHTML = ""; // Очищаем предыдущие карточки

                // Сортируем данные по UID
                data1.sort((a, b) => a.uid - b.uid);

                data1.forEach(vds => {
                    const cardHtml = `
                    <div class="col-md-4 mb-4">
                        <div class="card">
                            <div class="card-header ${vds.Status ? 'bg-success' : 'bg-danger'} text-white d-flex justify-content-between">
                                <span>${vds.Status ? 'On' : 'Off'}</span>
                            </div>
                            <div class="card-body">
                                <h5 class="card-title">Name: ${vds.Name}</h5>
                                <p class="card-text">IP/Port: ${vds.Ip}</p>
                                <div class="d-flex justify-content-between">
                                    <button class="btn btn-primary vds-btn" data-vds-id="${vds.uid}">
                                        <i class="bi bi-power"></i>
                                    </button>
                                    <button class="btn btn-danger delete-vds-btn" data-vds-id="${vds.uid}">
                                        <i class="bi bi-trash"></i>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>`;
                    cardsContainer.insertAdjacentHTML("beforeend", cardHtml);
                });

                // Находим все кнопки включения VDS и добавляем обработчики событий
                const vdsBtns = document.querySelectorAll(".vds-btn");
                vdsBtns.forEach(vdsBtn => {
                    vdsBtn.addEventListener("click", () => {
                        const vdsId = vdsBtn.dataset.vdsId;
                        this.turnVds(vdsId);
                    });
                });

                // Находим все кнопки удаления VDS и добавляем обработчики событий
                const deleteVdsBtns = document.querySelectorAll(".delete-vds-btn");
                deleteVdsBtns.forEach(deleteVdsBtn => {
                    deleteVdsBtn.addEventListener("click", () => {
                        const vdsId = deleteVdsBtn.dataset.vdsId;
                        this.deleteVds1(vdsId);
                    });
                });

            } else {
                console.error("Data is not an array:", data1);
            }
        }, error => {
            console.error("Error fetching data:", error);
        });
    };

    this.getVdsAdmin = () => {
        util.get("/getVdsTableAdmin", data1 => {
            console.log(data1, "Data from /getVdsTableAdmin");
            if (Array.isArray(data1)) {
                vdsTableAdmin.style.display = "block";
                util.id("addVdsBtnAdmin").style.display = "inline-block";
                util.id("deleteVdsBtn").style.display = "inline-block";
                this.vdsTableBodyAdmin = util.id("vdsTableAdmin").querySelector("tbody"); // Сохраняем ссылку на tbody
                this.vdsTableBodyAdmin.innerHTML = ""; // Очищаем текущие строки

                data1.forEach(vds => {
                    const row = `<tr>
                    <td>${vds.uid}</td>
                    <td>${vds.Uid}</td>
                    <td>${vds.Name}</td>
                    <td>${vds.Ip}</td>
                    <td>${vds.Login}</td>
                    <td>${vds.Pass}</td>
                    <td>${vds.Status ? 'Вкл' : 'Выкл'}</td>  
                </tr>`;
                    this.vdsTableBodyAdmin.innerHTML += row;
                });
            } else {
                console.error("Data is not an array:", data1);
            }
        }, error => {
            console.error("Error fetching data:", error);
        });
    };

    this.turnVds = (vdsId) => {
        console.log("Включение VDS id:", vdsId)
        util.post('/turnVds', {vid: parseInt(vdsId,10)}, (response) => {
            console.log("Response:", response)
            // Обработка ответа сервера
            if (response.success) {
                this.getVds();
            } else {
                // Ошибка при удалении VDS
                console.log("Ошибка при включении VDS id:", vdsId)
                alert(response.message);
            }
        });
    };

    this.deleteVds1 = (vdsId) => {
        console.log("Удаление VDS id:", vdsId)
        util.post('/deleteVds', {vid: parseInt(vdsId,10)}, (response) => {
            console.log("Response:", response)
            // Обработка ответа сервера
            if (response.success) {
                this.getVds();
            } else {
                // Ошибка при удалении VDS
                console.log("Ошибка при удалении VDS id:", vdsId)
                alert(response.message);
            }
        });
    };

    this.deleteVds = () => {
        // Вызываем модальное окно
        util.modal('confirmDeleteModalVds', 'show');

        // Получаем кнопку подтверждения удаления и добавляем обработчик событий
        document.getElementById('confirmDeleteBtnVds').onclick = () => {
            const vdsId = parseInt(document.getElementById('deleteVdsIdInput').value, 10);

            util.post('/deleteVds', {vid: vdsId}, (response) => {

                // Обработка ответа сервера
                if (response.success) {
                    if (util.id("addVdsBtnAdmin").style.display==="none") {
                        this.getVds();
                    } else{
                        this.getVdsAdmin();
                    }
                } else {
                    // Ошибка при удалении VDS
                    alert(response.message);
                }
            });

            // Закрыть модальное окно
            util.modal('confirmDeleteModalVds', 'hide');
        };
    };

    this.authIn = () => {
        //console.log("flskndlkfdsn11")
        util.post("/auth", {
            login: util.id("authLogin").value,
            pass: util.id("authPassword").value
        }, data => {
            //console.log(data,"qqq")
            if (data["token"]) {
                //console.log(data,"qqqw")
                localStorage.setItem("token", data["token"]);
                util.setToken(data["token"]);
                this.showUsersTable();
            } else {
                alert("Ytf");
            }
        }, error => {
            console.error("Error authenticating:", error);
        });
    };

    this.init();
}

document.getElementById('logoutBtn').addEventListener('click', function() {
    fetch('/logout', {
        method: 'POST',
        headers: {
            'Authorization': localStorage.getItem('token')
        }
    })
        .then(response => response.text())
        .then(data => {
            console.log(data);
            localStorage.removeItem('token');
            window.location.reload();
        })
        .catch(error => console.error('Ошибка:', error));
});

document.getElementById('deleteUserBtn').addEventListener('click', function() {
    main.deleteUser();
});

const main = new Main();