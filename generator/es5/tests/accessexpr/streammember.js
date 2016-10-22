$module('streammember', function () {
  var $static = this;
  this.$class('8933ac2f', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.SomeInt = $t.fastbox(2, $g.____testlib.basictypes.Integer);
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.AnotherThing = function (somestream) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.streamaccess(somestream, 'SomeInt');
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
