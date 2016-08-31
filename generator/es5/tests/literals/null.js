$module('null', function () {
  var $static = this;
  $static.DoSomething = function () {
    return $promise.empty();
  };
});
