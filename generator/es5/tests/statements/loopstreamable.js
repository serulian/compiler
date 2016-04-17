$module('loopstreamable', function () {
  var $static = this;
  this.$class('SomeStreamable', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Stream = function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.loopstreamable.SomeStream.new().then(function ($result0) {
                $result = $result0;
                $state.current = 1;
                $continue($state);
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

  this.$class('SomeStream', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.wasChecked = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Next = function () {
      var $this = this;
      var r;
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              r = $this.wasChecked;
              $this.wasChecked = $t.box(true, $g.____testlib.basictypes.Boolean);
              $g.____testlib.basictypes.Tuple($g.____testlib.basictypes.Boolean, $g.____testlib.basictypes.Boolean).Build($t.box(true, $g.____testlib.basictypes.Boolean), $t.box(!$t.unbox(r), $g.____testlib.basictypes.Boolean)).then(function ($result0) {
                $result = $result0;
                $state.current = 1;
                $continue($state);
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
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.box(1234, $g.____testlib.basictypes.Integer);
            $state.current = 1;
            continue;

          case 1:
            somethingElse.Stream().then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $temp1 = $result;
            $state.current = 3;
            continue;

          case 3:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $state.current = 4;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
            $result;
            something = $temp0.First;
            if ($t.unbox($temp0.Second)) {
              $state.current = 5;
              continue;
            } else {
              $state.current = 6;
              continue;
            }
            break;

          case 5:
            $t.box(7654, $g.____testlib.basictypes.Integer);
            $state.current = 3;
            continue;

          case 6:
            $t.box(5678, $g.____testlib.basictypes.Integer);
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
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            result = $t.box('noloop', $g.____testlib.basictypes.String);
            $g.loopstreamable.SomeStreamable.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            s = $result;
            $state.current = 2;
            continue;

          case 2:
            s.Stream().then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $temp1 = $result;
            $state.current = 4;
            continue;

          case 4:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $state.current = 5;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 5:
            $result;
            i = $temp0.First;
            if ($t.unbox($temp0.Second)) {
              $state.current = 6;
              continue;
            } else {
              $state.current = 7;
              continue;
            }
            break;

          case 6:
            result = i;
            $state.current = 4;
            continue;

          case 7:
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
