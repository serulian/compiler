$module('null', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      null;
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
