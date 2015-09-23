angular.module('uiRouterSample.cars.service', ['ngCookies'])


.factory('carsFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var carsFactory = {};

    carsFactory.all = function() {
        var url = 'http://crunchy-pg-client:13002/car/list';
        console.log(url);

        return $http.get(url);
    };


    carsFactory.get = function(carid) {

        var url = 'http://crunchy-pg-client:13002/car/' + carid;
        console.log(url);

        return $http.get(url);
    };


    carsFactory.delete = function(carid) {

        var url = 'http://crunchy-pg-client:13002/car/delete';
        console.log(url);

        return $http.post(url, {
		'ID' : carid,
	});
    };

    carsFactory.update = function(car) {

        var url = 'http://crunchy-pg-client:13002/car/update';
        console.log(url);

        return $http.post(url, {
		'ID' : car.ID,
		'Model' : car.Model,
		'Brand' : car.Brand,
		'Year' : car.Year,
		'Price' : car.Price,
	});
    };


    carsFactory.add = function(car) {


        var url = 'http://crunchy-pg-client:13002' + '/car/add';
        console.log(url);

        return $http.post(url, {
		'Model' : car.Model,
		'Brand' : car.Brand,
		'Year' : car.Year,
		'Price' : car.Price,
	});
    };

    return carsFactory;
}]);
