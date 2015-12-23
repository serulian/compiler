$module('conditional', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              if (true) {
                $state.current = 1;
              } else {
                $state.current = 2;
              }
              continue;

            case 1:
              123;
              $state.current = 2;
              continue;

            default:
              $state.current = -1;
              return;
          }
        }
      } catch (e) {
        $state.error = e;
        $state.current = -1;
        $callback($state);
      }
    };
    return $promise.build($state);
  };
});
