$module('loopvar', function () {
  var $static = this;
  $static.DoSomething = function (somethingElse) {
    var something;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            1234;
            $state.current = 1;
            continue;

          case 1:
            somethingElse.Next(function ($hasNext, $nextItem) {
              if ($hasNext) {
                something = $nextItem;
                $state.current = 2;
              } else {
                $state.current = 3;
              }
              $state.$next($callback);
            });
            return;

          case 2:
            7654;
            $state.current = 1;
            continue;

          case 3:
            5678;
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
