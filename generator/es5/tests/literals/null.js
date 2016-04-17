$module('null', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($continue) {
      null;
      $state.resolve();
    });
    return $promise.build($state);
  };
});
