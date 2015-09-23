angular.module('uiRouterSample.home.service', ['ngCookies'])

// A RESTful factory for retrieving home from 'home.json'
.factory('homeFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, utils) {

    var factory = {};

    factory.getconn = function() {
     	var url = 'http://crunchy-pg-client:13002/conn/list';
       	console.log(url);
	return $http.get(url);
    };

    return factory;
}]);
