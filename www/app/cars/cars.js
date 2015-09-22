angular.module('uiRouterSample.cars', [
    'ui.router',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('cars', {

                    abstract: true,

                    url: '/cars',

                    templateUrl: 'app/cars/cars.html',

                    resolve: {
                        cars: ['$cookieStore', 'carsFactory',
                            function($cookieStore, carsFactory) {
                                if (!$cookieStore.get('cpm_token')) {
                                    var nothing = [];
                                    console.log('returning nothing');
                                    return nothing;
                                }

                                return carsFactory.all();
                            }
                        ]
                    },

                    controller: ['$scope', '$state', '$cookieStore', 'utils', 'cars',
                        function($scope, $state, $cookieStore, utils, cars) {

                            if (!$cookieStore.get('cpm_token')) {
                                $state.go('login', {
                                    userId: 'hi'
                                });
                            }

                            $scope.cars = cars;

                            $scope.goToFirst = function() {
                                if ($scope.cars.data.length > 0) {
                                    //console.log($scope.cars.data[0].ID);
                                    var randId = $scope.cars.data[0].ID;

                                    $state.go('cars.detail.details', {
                                        serverId: randId
                                    });
                                }
                            };
                            $scope.goToFirst();
                        }
                    ]
                })

            .state('cars.list', {

                url: '',

                templateUrl: 'app/cars/cars.list.html',
                controller: ['$scope', '$state', '$stateParams', 'carsFactory', 'cars', 'utils',
                    function($scope, $state, $stateParams, carsFactory, cars, utils) {

                        carsFactory.all()
                            .success(function(data) {
                                //console.log('successful get in list =' + JSON.stringify(data));
                                $scope.cars = data;
                            })
                            .error(function(error) {
                                $scope.alerts = [{
                                    type: 'danger',
                                    msg: error.message
                                }];
                                console.log('here is an error ' + error.message);
                            });


                        if ($scope.cars.data.length > 0) {
                        	var randId = $scope.cars.data[0].ID;
                            $state.go('cars.detail.details', {
                                serverId: randId
                            });
                        }
                    }
                ]

            })


            .state('cars.detail', {

                url: '/{serverId:[0-9]{1,4}}',

                views: {

                    '': {
                        templateUrl: 'app/cars/cars.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'carsFactory', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, carsFactory, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in cars');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                if ($scope.cars.data.length > 0) {
                                    angular.forEach($scope.cars.data, function(item) {
                                        if (item.ID == $stateParams.serverId) {
                                            $scope.server = item;
                                        }
                                    });
                                }

                            }
                        ]
                    },

                }
            })

            .state('cars.detail.item', {
                url: '/item/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/cars/cars.detail.item.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.server.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('cars.detail.add', {
                url: '/add/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/cars/cars.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'carsFactory', 'utils',
                            function($scope, $stateParams, $state, carsFactory, utils) {

                                var newserver = {};
                                newserver.ID = '0';
                                newserver.Name = 'newserver';
                                newserver.IPAddress = '1.1.1.1';
                                newserver.DockerBridgeIP = '172.17.42.1';
                                newserver.PGDataPath = '/var/cpm/data/pgsql';
                                newserver.ServerClass = 'low';
                                newserver.CreateDate = '00';
                                $scope.server = newserver;

                                $scope.add = function() {
                                    $scope.server.ID = 0; //0 means to do an insert

                                    carsFactory.add($scope.server)
                                        .success(function(data) {
                                            $state.go('cars.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                        });
                                };
                            }
                        ]
                    },
                }
            })

            .state('cars.detail.details', {
                url: '/details/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/cars/cars.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'carsFactory', 'utils',
                            function($scope, $stateParams, $state, carsFactory, utils) {
                                $scope.save = function() {
                                    carsFactory.add($scope.server)
                                        .success(function(data) {
                                            $scope.alerts = [{
                                                type: 'success',
                                                msg: 'saved'
                                            }];
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                            }
                        ]
                    },
                }
            })

            .state('cars.detail.delete', {
                url: '/delete/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/cars/cars.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'carsFactory', 'utils',
                            function($scope, $stateParams, $state, carsFactory, utils) {

                                var server = $scope.server;
                                $scope.delete = function() {
                                    carsFactory.delete($stateParams.serverId)
                                        .success(function(data) {
                                            $state.go('cars.list', {});
                                            $state.go('cars.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                            }
                        ]
                    },
                }
            })




        }
    ]
);
