$module('loopvar', function () {
  var $static = this;
  this.$class('SomeStream', false, function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve(false).then(function (result) {
        instance.wasChecked = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Next = function () {
      var $this = this;
      var r;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              r = $this.wasChecked;
              $this.wasChecked = true;
              $g.____testlib.basictypes.Tuple($g.____testlib.basictypes.Boolean, $g.____testlib.basictypes.Boolean).Build(true, !r).then(function ($result0) {
                $result = $result0;
                $state.current = 1;
                $callback($state);
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
              $state.resolve($result);
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

  $static.DoSomething = function (somethingElse) {
    var $temp0;
    var $temp1;
    var something;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            1234;
            $state.current = 1;
            continue;

          case 1:
            $temp1 = somethingElse;
            $state.current = 2;
            continue;

          case 2:
            $temp1.Next().then(function ($result0) {
              $result = $temp0 = $result0;
              $state.current = 3;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $result;
            something = $temp0.First;
            if ($temp0.Second) {
              $state.current = 4;
              continue;
            } else {
              $state.current = 5;
              continue;
            }
            break;

          case 4:
            7654;
            $state.current = 2;
            continue;

          case 5:
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
  $static.TEST = function () {
    var $temp0;
    var $temp1;
    var i;
    var result;
    var s;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            result = '';
            $g.loopvar.SomeStream.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            s = $result;
            $state.current = 2;
            continue;

          case 2:
            $temp1 = s;
            $state.current = 3;
            continue;

          case 3:
            $temp1.Next().then(function ($result0) {
              $result = $temp0 = $result0;
              $state.current = 4;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
            $result;
            i = $temp0.First;
            if ($temp0.Second) {
              $state.current = 5;
              continue;
            } else {
              $state.current = 6;
              continue;
            }
            break;

          case 5:
            result = i;
            $state.current = 3;
            continue;

          case 6:
            $state.resolve(result);
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
