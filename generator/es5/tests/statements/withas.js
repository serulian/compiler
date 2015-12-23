$module('withas', function () {
  var $static = this;
  $static.DoSomething = function (someExpr) {
    var someName;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            123;
            someName = someExpr;
            $t.pushr($state, 'someName', someName);
            456;
            $t.popr($state, 'someName');
            789;
            $state.current = -1;
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
