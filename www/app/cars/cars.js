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
                                return carsFactory.all();
                            }
                        ]
                    },

                    controller: ['$scope', '$state', '$cookieStore', 'utils', 'cars',
                        function($scope, $state, $cookieStore, utils, cars) {

                            $scope.cars = cars;

                            $scope.goToFirst = function() {
			    	console.log(JSON.stringify($scope.cars));
                                if ($scope.cars.data.length > 0) {
                                    //console.log($scope.cars.data[0].ID);
                                    var randId = $scope.cars.data[0].ID;

                                    $state.go('cars.detail.details', {
                                        carId: randId
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
                                carId: randId
                            });
                        }
                    }
                ]

            })


            .state('cars.detail', {

                url: '/{carId:[0-9]{1,4}}',

                views: {

                    '': {
                        templateUrl: 'app/cars/cars.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'carsFactory', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, carsFactory, utils) {

                                if ($scope.cars.data.length > 0) {
                                    angular.forEach($scope.cars.data, function(item) {
                                        if (item.ID == $stateParams.carId) {
                                            $scope.car = item;
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
                                $scope.item = utils.findById($scope.car.items, $stateParams.itemId);

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

                                var newcar = {};
                                newcar.ID = '0';
                                newcar.Name = 'newcar';
                                newcar.IPAddress = '1.1.1.1';
                                newcar.DockerBridgeIP = '172.17.42.1';
                                newcar.PGDataPath = '/var/cpm/data/pgsql';
                                newcar.ServerClass = 'low';
                                newcar.CreateDate = '00';
                                $scope.car = newcar;

                                $scope.add = function() {
                                    $scope.car.ID = 0; //0 means to do an insert

                                    carsFactory.add($scope.car)
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
                                    carsFactory.update($scope.car)
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

                                var car = $scope.car;
                                $scope.delete = function() {
                                    carsFactory.delete($stateParams.carId)
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
