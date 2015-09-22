angular.module('uiRouterSample.cars.service', ['ngCookies'])


.factory('carsFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var carsFactory = {};

    carsFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/cars/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    carsFactory.get = function(serverid) {

        var url = $cookieStore.get('AdminURL') + '/server/' + serverid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    carsFactory.delete = function(serverid) {

        var url = $cookieStore.get('AdminURL') + '/deleteserver/' + serverid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    carsFactory.add = function(server) {

        var cleanip = server.IPAddress.replace(/\./g, "_");
        var cleandockerip = server.DockerBridgeIP.replace(/\./g, "_");
        var cleanpath = server.PGDataPath.replace(/\//g, "_");

        var url = $cookieStore.get('AdminURL') + '/addserver/' +
            server.ID + '.' +
            server.Name + '.' +
            cleanip + '.' +
            cleandockerip + '.' +
            cleanpath + '.' +
            server.ServerClass + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    return carsFactory;
}]);
